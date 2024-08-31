/*
Copyright 2024 asjfoajs.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dappsv1 "github.com/asjfoajs/MyOperatorProjects/application-operator/api/v1"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.hyj.cn,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.hyj.cn,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.hyj.cn,resources=applications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	//get the Application
	//声明一个*Application类型的实例app用来接收我们的CR
	app := &dappsv1.Application{}
	//NamespacedName在这里也就是default/application-sample
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		//err分很多中情况，如果找不到，一般不需要进一步处理，只是说明这个CR被删了而已
		if errors.IsNotFound(err) {
			l.Info("Application resource not found. Ignoring since object must be deleted")
			//直接返回，不带错误，结束本次调谐
			return ctrl.Result{}, nil
		}
		//除了NotFound之外的错误，比如连不上apiserver等，这时需要打印错误信息，然后返回这个
		//错误以及表示1分钟后重试的Result
		l.Error(err, "Failed to get Application")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, err
	}

	//create pods
	for i := 0; i < int(app.Spec.Replicas); i++ {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Name + "-" + strconv.Itoa(i),
				Namespace: app.Namespace,
				Labels:    app.Labels,
			},
			Spec: app.Spec.Teplate.Spec,
		}

		if err := r.Create(ctx, pod); err != nil {
			l.Error(err, "Failed to create pod")
			return ctrl.Result{RequeueAfter: 1 * time.Minute}, err
		}

		l.Info(fmt.Sprintf("the Pod (%s) has created", pod.Name))
	}

	l.Info("all pods has created")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dappsv1.Application{}).
		Complete(r)
}
