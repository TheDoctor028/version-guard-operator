package controller

import (
	"context"
	"github.com/TheDoctor028/version-guard-operator/internal/model"
	"time"

	versionguradv1alpha1 "github.com/TheDoctor028/version-guard-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	Notifier model.Notifier
}

// +kubebuilder:rbac:groups=version-gurad.tmit.bme.hu,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=version-gurad.tmit.bme.hu,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=version-gurad.tmit.bme.hu,resources=applications/finalizers,verbs=update

// Reconcile is the main controller loop for Application resources
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	app := &versionguradv1alpha1.Application{}
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.Notifier.SendChangeNotification(model.VersionChangeData{
		Kind:          "Application",
		Name:          app.Name,
		Selector:      model.ParseSelector(app.Spec.Selector),
		ContainerName: app.Spec.Name,
		Namespace:     app.Namespace,
		Image:         app.Spec.Image,
		Timestamp:     time.Now().UTC(),
	}); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&versionguradv1alpha1.Application{}).
		Complete(r)
}
