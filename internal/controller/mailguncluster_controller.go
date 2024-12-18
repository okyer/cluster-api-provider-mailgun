/*
Copyright 2024.

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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1alpha1 "github.com/okyer/cluster-api-provider-mailgun/api/v1alpha1"
	"github.com/pkg/errors"
)

// MailgunClusterReconciler reconciles a MailgunCluster object
type MailgunClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	// Mailgun   mailgun.Mailgun
	Recipient string
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=mailgunclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=mailgunclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=mailgunclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MailgunCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *MailgunClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logr := log.FromContext(ctx)
	logr.WithValues("mailguncluster", req.NamespacedName)

	// TODO(user): your logic here
	var cluster infrastructurev1alpha1.MailgunCluster
	if err := r.Get(ctx, req.NamespacedName, &cluster); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if cluster.Status.MessageID != nil {
		// We already sent a message, so skip reconciliation
		return ctrl.Result{}, nil
	}

	// subject := fmt.Sprintf("[%s] New Cluster %s requested", cluster.Spec.Priority, cluster.Name)
	// body := fmt.Sprintf("Hello! One cluster please.\n\n%s\n", cluster.Spec.Request)

	// msg := mailgun.NewMessage(cluster.Spec.Requester, subject, body, r.Recipient)
	// _, msgID, err := r.Mailgun.Send(msg)
	// if err != nil {
	// 	return ctrl.Result{}, err
	// }

	msgID := "test"

	helper, err := patch.NewHelper(&cluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	cluster.Status.MessageID = &msgID
	if err := helper.Patch(ctx, &cluster); err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "couldn't patch cluster %q", cluster.Name)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MailgunClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.MailgunCluster{}).
		Named("mailguncluster").
		Complete(r)
}
