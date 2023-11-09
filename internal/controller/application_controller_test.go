package controller

import (
	"context"
	"fmt"
	"github.com/TheDoctor028/version-guard-operator/api/v1alpha1"
	"github.com/TheDoctor028/version-guard-operator/internal/mocks/notifier_mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Application controller", func() {
	It("should send the data through Notifier if a new Application CR is created", func() {
		reconciler, mockCtrl, notifierMock := setupApplicationReconciler()

		app := createTestApplication(createApplicationTemplate())

		notifierMock.EXPECT().SendChangeNotification(gomock.Any()).Return(nil)

		res, err := reconciler.Reconcile(context.TODO(), ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: app.Namespace,
				Name:      app.Name,
			},
		})
		assert.Nil(GinkgoT(), err)
		assert.NotNil(GinkgoT(), res)

		mockCtrl.Finish()
	})

	It("should send the data through Notifier if a new Application CR is created and again when updated", func() {
		reconciler, mockCtrl, notifierMock := setupApplicationReconciler()

		app := createTestApplication(createApplicationTemplate())
		defer k8sClient.Delete(context.TODO(), app)

		notifierMock.EXPECT().SendChangeNotification(gomock.Any()).Return(nil)
		res, err := reconciler.Reconcile(context.TODO(), ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: app.Namespace,
				Name:      app.Name,
			},
		})
		assert.Nil(GinkgoT(), err)
		assert.NotNil(GinkgoT(), res)

		app.Spec.Image = "nginx:1.19.2"
		err = k8sClient.Update(context.TODO(), app)
		assert.Nil(GinkgoT(), err)

		notifierMock.EXPECT().SendChangeNotification(gomock.Any()).Return(nil)
		res, err = reconciler.Reconcile(context.TODO(), ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: app.Namespace,
				Name:      app.Name,
			},
		})
		assert.Nil(GinkgoT(), err)
		assert.Equal(GinkgoT(), reconcile.Result{Requeue: false, RequeueAfter: 0}, res)

		mockCtrl.Finish()
	})

	It("should return an error if the SendChangeNotification fails", func() {
		reconciler, mockCtrl, notifierMock := setupApplicationReconciler()

		app := createTestApplication(createApplicationTemplate())
		defer k8sClient.Delete(context.TODO(), app)

		notifierMock.EXPECT().SendChangeNotification(gomock.Any()).Return(assert.AnError)
		res, err := reconciler.Reconcile(context.TODO(), ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: app.Namespace,
				Name:      app.Name,
			},
		})
		assert.Error(GinkgoT(), err, assert.AnError)
		assert.Equal(GinkgoT(), reconcile.Result{Requeue: false, RequeueAfter: 0}, res)

		mockCtrl.Finish()

	})

	It("should return an nil error if the Application CR docent exists", func() {
		reconciler, mockCtrl, _ := setupApplicationReconciler()

		res, err := reconciler.Reconcile(context.TODO(), ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "not-existing-app",
			},
		})
		assert.Nil(GinkgoT(), err)
		assert.Equal(GinkgoT(), reconcile.Result{Requeue: false, RequeueAfter: 0}, res)

		mockCtrl.Finish()

	})
})

func setupApplicationReconciler() (ApplicationReconciler, *gomock.Controller, *notifier_mock.MockNotifier) {
	mockCtr := gomock.NewController(GinkgoT())
	mockNotifier := notifier_mock.NewMockNotifier(mockCtr)
	reconciler := ApplicationReconciler{
		Client: k8sClient,
		Scheme: k8sClient.Scheme(),

		Notifier: mockNotifier,
	}

	return reconciler, mockCtr, mockNotifier
}

func createApplicationTemplate() *v1alpha1.Application {
	return &v1alpha1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("test-app-%s", rand.String(5)),
			Namespace: "default",
		},
		Spec: v1alpha1.ApplicationSpec{
			Name:  "nginx",
			Image: "nginx:latest",
			Selector: map[string]string{
				"app": "nginx",
			},
		},
	}
}

func createTestApplication(app *v1alpha1.Application) *v1alpha1.Application {
	err := k8sClient.Create(context.TODO(), app)
	assert.Nil(GinkgoT(), err)

	return app
}
