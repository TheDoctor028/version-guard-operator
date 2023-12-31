package controller

import (
	"context"
	"fmt"
	"github.com/TheDoctor028/version-guard-operator/internal/mocks/notifier_mock"
	"github.com/TheDoctor028/version-guard-operator/internal/model"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

var _ = Describe("Deployment controller", func() {
	When("a new Deployment is created", func() {
		It("should not send Notification if there is no Application for the Deployment", func() {
			reconciler, mockCtrl, notifierMock := setupDeploymentReconciler()

			deployment := createDeploymentTemplate()
			deployment = createTestDeployment(deployment)

			notifierMock.EXPECT().SendChangeNotification(gomock.Any()).Return(nil).Times(0)
			res, err := reconciler.Reconcile(context.TODO(), ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: deployment.Namespace,
					Name:      deployment.Name,
				},
			})

			assert.Nil(GinkgoT(), err)
			assert.NotNil(GinkgoT(), res)
			mockCtrl.Finish()
		})
		It("should send Notification if there is Application for the Deployment", func() {
			reconciler, mockCtrl, notifierMock := setupDeploymentReconciler()

			deployment := createDeploymentTemplate()
			deployment.SetLabels(map[string]string{
				"app": "nginx",
			})
			deployment = createTestDeployment(deployment)

			_ = createTestApplication(createApplicationTemplate())

			notifierMock.EXPECT().SendChangeNotification(model.VersionChangeDataEQ(model.VersionChangeData{
				Kind:          model.DeploymentKind,
				Name:          deployment.Name,
				Namespace:     deployment.Namespace,
				Selector:      model.ParseSelector(deployment.Spec.Selector.MatchLabels),
				ContainerName: deployment.Spec.Template.Spec.Containers[0].Name,
				Image:         deployment.Spec.Template.Spec.Containers[0].Image,
				Timestamp:     time.Now(),
			})).Return(nil).Times(1)
			res, err := reconciler.Reconcile(context.TODO(), ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: deployment.Namespace,
					Name:      deployment.Name,
				},
			})

			assert.Nil(GinkgoT(), err)
			assert.NotNil(GinkgoT(), res)
			mockCtrl.Finish()
		})
	})
})

func setupDeploymentReconciler() (DeploymentReconciler, *gomock.Controller, *notifier_mock.MockNotifier) {
	mockCtr := gomock.NewController(GinkgoT())
	mockNotifier := notifier_mock.NewMockNotifier(mockCtr)
	reconciler := DeploymentReconciler{
		Client: k8sClient,
		Scheme: k8sClient.Scheme(),

		Notifier: mockNotifier,
	}

	return reconciler, mockCtr, mockNotifier
}

func createDeploymentTemplate() *v1.Deployment {
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("test-deployment-%s", rand.String(5)),
			Namespace: "default",
		},
		Spec: v1.DeploymentSpec{
			Replicas: nil,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
						},
					},
				},
			},
		},
	}
}

func createTestDeployment(deployment *v1.Deployment) *v1.Deployment {
	err := k8sClient.Create(context.TODO(), deployment)
	assert.Nil(GinkgoT(), err)

	return deployment
}
