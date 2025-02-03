/*
Copyright 2025.

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
	"io"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/example/tutorial-gitops-operator/api/v1alpha1"
	gitopsv1alpha1 "github.com/example/tutorial-gitops-operator/api/v1alpha1"
)

// O2ImsReconciler reconciles a O2Ims object
type O2ImsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gitops.example.com,resources=o2ims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gitops.example.com,resources=o2ims/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gitops.example.com,resources=o2ims/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the O2Ims object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *O2ImsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here

	url_base := ""
	o2ims := &v1alpha1.O2Ims{}
	err := r.Get(ctx, req.NamespacedName, o2ims)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("memcached resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get memcached")
		return ctrl.Result{}, err
	}

	o2ims.Status.Phase = "Finished"
	url := url_base + o2ims.Spec.Endpoint
	if o2ims.Spec.DeploymentManager != "" {
		url = url + o2ims.Spec.DeploymentManager
	}
	log.Info(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err, "HTTP request failed")
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err, "HTTP request failed")
		} else {
			o2ims.Status.RetrievedInformation = string(body)
			log.Info(string(body))
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *O2ImsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gitopsv1alpha1.O2Ims{}).
		Complete(r)
}
