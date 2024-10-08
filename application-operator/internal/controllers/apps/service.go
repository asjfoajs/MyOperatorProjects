package controller

import (
	"context"
	v1 "github.com/asjfoajs/MyOperatorProjects/application-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *ApplicationReconciler) reconcileService(ctx context.Context, app *v1.Application) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var svc = &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      app.Name,
		Namespace: app.Namespace,
	}, svc)

	if err == nil {
		log.Info("The Service has already exist.")
		if reflect.DeepEqual(svc.Status, app.Status.Network) {
			return ctrl.Result{}, nil
		}

		app.Status.Network = svc.Status
		if err := r.Status().Update(ctx, app); err != nil {
			log.Error(err, "Failed to update Application status")
			return ctrl.Result{RequeueAfter: GenericRequeueDuration}, err
		}
		log.Info("The application status has been update")

		return ctrl.Result{}, nil
	}

	if !errors.IsNotFound(err) {
		log.Error(err, "Failed to get Service,will requeue after a short time")
		return ctrl.Result{RequeueAfter: GenericRequeueDuration}, err
	}

	newSvc := &corev1.Service{}
	newSvc.SetName(app.Name)
	newSvc.SetNamespace(app.Namespace)
	newSvc.SetLabels(app.Labels)
	newSvc.Spec = app.Spec.Service.ServiceSpec
	newSvc.Spec.Selector = app.Labels

	if err := ctrl.SetControllerReference(app, newSvc, r.Scheme); err != nil {
		log.Error(err, "Failed to set controller reference , will requeue after a short time.")
		return ctrl.Result{RequeueAfter: GenericRequeueDuration}, err
	}

	if err := r.Create(ctx, newSvc); err != nil {
		log.Error(err, "Failed to create Service,will requeue after a short time")
		return ctrl.Result{RequeueAfter: GenericRequeueDuration}, err
	}

	log.Info("The Service has been created.")
	return ctrl.Result{}, nil
}
