package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/TheDoctor028/version-guard-operator/internal/model"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

const containerVerAnnotation = "container-img-version.version-guard.tmit.bme.hu/%s"

var ignoredNamespaces = []string{"kube-system", "kube-public", "kube-node-lease"}

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Ignore deployments in the ignored namespaces
	for _, namespace := range ignoredNamespaces {
		if req.Namespace == namespace {
			return reconcile.Result{}, nil
		}
	}

	deployment := &appsv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, deployment); err != nil {
		return reconcile.Result{}, err
	}

	for _, container := range deployment.Spec.Template.Spec.Containers {
		if hasImageChanged(deployment, &container) {
			if err := r.sendChangeToAPI(deployment, &container); err != nil {
				return reconcile.Result{}, err
			}

			if err := r.addAnnotationsContainerVerToDeployment(ctx, deployment, container); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	return reconcile.Result{}, nil

}

// addAnnotationsContainerVerToDeployment adds an annotations to the deployment with the current image version of the container
func (r *DeploymentReconciler) addAnnotationsContainerVerToDeployment(ctx context.Context, deployment *appsv1.Deployment, container v1.Container) error {
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

func (r *DeploymentReconciler) sendChangeToAPI(deployment *appsv1.Deployment, container *v1.Container) error {
	data := model.DeploymentData{
		Kind:          "Deployment",
		Name:          deployment.Name,
		Namespace:     deployment.Namespace,
		ContainerName: container.Name,
		Image:         container.Image,
		Timestamp:     time.Now().UTC(),
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		return fmt.Errorf("API_URL environment variable is not set")
	}

	// Call the API endpoint with the JSON data
	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API call failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}