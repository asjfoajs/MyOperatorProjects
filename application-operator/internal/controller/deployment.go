package controller

import (
	"context"
	v1 "github.com/asjfoajs/MyOperatorProjects/application-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *ApplicationReconciler) reconcileDeployment(ctx context.Context, app *v1.Application) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var dp = &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      app.Name,
		Namespace: app.Namespace,
	}, dp)

	if err == nil {
		log.Info("The Deployment has already exist.")
		if reflect.DeepEqual(dp.Status, app.Status.Workflow) {
			return ctrl.Result{}, nil
		}

		app.Status.Workflow = dp.Status
		if err := r.Status().Update(ctx, app); err != nil {
			log.Error(err, "Failed to update Application status")
			return ctrl.Result{RequeueAfter: GenericRequeueDuration}, err
		}

		log.Info("The Deployment has been update.")
		return ctrl.Result{}, nil
	}

	if !errors.IsNotFound(err) {
		log.Error(err, "Failed to get Deployment,will requeue after a short time")
		return ctrl.Result{RequeueAfter: GenericRequeueDuration}, err
	}

	newDp := &appsv1.Deployment{}
	newDp.SetName(app.Name)
	newDp.SetNamespace(app.Namespace)
	newDp.SetLabels(app.Labels)
	newDp.Spec = app.Spec.Deployment.DeploymentSpec
	newDp.Spec.Template.SetLabels(app.Labels)

	if err := ctrl.SetControllerReference(app, newDp, r.Scheme); err != nil {
		log.Error(err, "Failed to set controller reference , will requeue after a short time.")
		return ctrl.Result{RequeueAfter: GenericRequeueDuration}, err
	}

	if err := r.Create(ctx, newDp); err != nil {
		log.Error(err, "Failed to create Deployment, will requeue after a short time.")
		return ctrl.Result{RequeueAfter: GenericRequeueDuration}, err
	}

	log.Info("The Deployment has been created.")
	return ctrl.Result{}, nil
}
