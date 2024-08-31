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

package v1

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var applicationlog = logf.Log.WithName("application-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *Application) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-apps-hyj-cn-v1-application,mutating=true,failurePolicy=fail,sideEffects=None,groups=apps.hyj.cn,resources=applications,verbs=create;update,versions=v1,name=mapplication.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Application{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Application) Default() {
	applicationlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.

	if r.Spec.Deployment.Replicas == nil {
		r.Spec.Deployment.Replicas = new(int32)
		*r.Spec.Deployment.Replicas = 3
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-apps-hyj-cn-v1-application,mutating=false,failurePolicy=fail,sideEffects=None,groups=apps.hyj.cn,resources=applications,verbs=create;update,versions=v1,name=vapplication.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Application{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Application) ValidateCreate() (admission.Warnings, error) {
	applicationlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, r.validateApplication()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Application) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	applicationlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, r.validateApplication()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Application) ValidateDelete() (admission.Warnings, error) {
	applicationlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

func (r *Application) validateApplication() error {
	if *r.Spec.Deployment.Replicas > 10 {
		return fmt.Errorf("replicas tooo many error")
	}
	return nil
}
