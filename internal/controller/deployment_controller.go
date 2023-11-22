package controller

import (
	"context"
	"fmt"
	"github.com/TheDoctor028/version-guard-operator/api/v1alpha1"
	"github.com/TheDoctor028/version-guard-operator/internal/model"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

const containerVerAnnotation = "container-img-version.version-guard.tmit.bme.hu/%s"

var defaultIgnoredNamespaces = []string{"kube-system", "kube-public", "kube-node-lease"}

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	Notifier model.Notifier
}

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	deployment := &appsv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, deployment); client.IgnoreNotFound(err) != nil {
		return reconcile.Result{}, err
	}

	appList := &v1alpha1.ApplicationList{}
	err := r.List(ctx, appList, &client.ListOptions{Namespace: ""})
	if err != nil {
		return reconcile.Result{}, err
	}

	for _, app := range appList.Items {
		if app.Namespace != req.Namespace {
			return reconcile.Result{}, nil
		}
	}

	app := findApplicationForDeployment(appList, deployment)
	if app == nil {
		return reconcile.Result{}, nil
	}

	for _, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == app.Spec.Name && hasImageChanged(deployment, &container) {
			if err := r.Notifier.SendChangeNotification(model.VersionChangeData{
				Kind:          "Deployment",
				Name:          deployment.Name,
				Namespace:     deployment.Namespace,
				Selector:      model.ParseSelector(app.Spec.Selector),
				ContainerName: container.Name,
				Image:         container.Image,
				Timestamp:     time.Now().UTC(),
			}); err != nil {
				return reconcile.Result{}, err
			}

			if err := r.addAnnotationsContainerVerToDeployment(ctx, deployment, container); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	return reconcile.Result{}, nil
}

// findApplicationForDeployment finds the application that the deployment belongs to
func findApplicationForDeployment(appList *v1alpha1.ApplicationList, deployment *appsv1.Deployment) *v1alpha1.Application {
	for _, app := range appList.Items {
		for key, val := range app.Spec.Selector {
			if deployment.Labels[key] == val {
				return &app
			}
		}
	}
	return nil
}

// addAnnotationsContainerVerToDeployment adds an annotations to the deployment with the current image version of the container
func (r *DeploymentReconciler) addAnnotationsContainerVerToDeployment(ctx context.Context, deployment *appsv1.Deployment, container v1.Container) error {
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
	}, deployment); client.IgnoreNotFound(err) != nil {
		return err
	}
	if deployment.Annotations == nil {
		deployment.Annotations = make(map[string]string)
	}
	deployment.Annotations[fmt.Sprintf(containerVerAnnotation, container.Name)] = container.Image
	return r.Update(ctx, deployment)
}

// hasImageChanged checks if the image of the container has changed
// since the last reconciliation
func hasImageChanged(deployment *appsv1.Deployment, container *v1.Container) bool {
	annotationKey := fmt.Sprintf(containerVerAnnotation, container.Name)
	annotationVal, success := deployment.Annotations[annotationKey]
	return !success || annotationVal != container.Image
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
