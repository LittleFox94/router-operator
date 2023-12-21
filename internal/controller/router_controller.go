/*
Copyright 2023.

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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"praios.lf-net.org/littlefox/router-operator/internal/resources"
	"praios.lf-net.org/littlefox/router-operator/internal/resources/bird"

	routingv1alpha1 "praios.lf-net.org/littlefox/router-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RouterReconciler reconciles a Router object
type RouterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *RouterReconciler) reconciliationForRequest(ctx context.Context, req ctrl.Request) (*resources.RouterReconciliation, error) {
	ret := resources.RouterReconciliation{
		Router:   &routingv1alpha1.Router{},
		Sessions: make([]*routingv1alpha1.Session, 0),
		Peers:    make([]*routingv1alpha1.Peer, 0),
	}

	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, ret.Router); err != nil {
		return nil, fmt.Errorf("error fetching routingv1alpha1.Router: %w", err)
	}

	allSessions := routingv1alpha1.SessionList{}
	if err := r.Client.List(ctx, &allSessions); err != nil {
		return nil, fmt.Errorf("error listing sessions for routingv1alpha1.Router: %w", err)
	}

	for _, session := range allSessions.Items {
		if session.GetNamespace() == ret.Router.GetNamespace() && session.Spec.Router.Name == ret.Router.GetName() {
			peer := routingv1alpha1.Peer{}
			if err := r.Client.Get(ctx, client.ObjectKey{Namespace: session.GetNamespace(), Name: session.Spec.Peer.Name}, &peer); err != nil {
				return nil, fmt.Errorf(
					"error retrieving peer %q for routingv1alpha1.Session %q: %w",
					session.Spec.Peer.Name,
					session.GetName(),
					err,
				)
			}

			ret.Sessions = append(ret.Sessions, session.DeepCopy())
			ret.Peers = append(ret.Peers, &peer)
		}
	}

	return &ret, nil
}

func (r *RouterReconciler) updateStatus(ctx context.Context, reconciliation *resources.RouterReconciliation) error {
	reconciliation.Router.Status.ClusterSessions = uint(len(reconciliation.Sessions))
	reconciliation.Router.Status.ClusterAnnouncements = uint(len(reconciliation.Announcements))
	reconciliation.Router.Status.LastUpdateTime = metav1.Now()

	if reconciliation.Router.Status.Conditions == nil {
		reconciliation.Router.Status.Conditions = make([]routingv1alpha1.RouterCondition, 0)
	}

	if err := r.Client.Status().Update(ctx, reconciliation.Router); err != nil {
		return fmt.Errorf("error updating status of routingv1alpha1.Router: %w", err)
	}

	return nil
}

//+kubebuilder:rbac:groups=routing.lf-net.org,resources=routers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=routing.lf-net.org,resources=routers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=routing.lf-net.org,resources=routers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Router object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *RouterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reconciliation, err := r.reconciliationForRequest(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := r.updateStatus(ctx, reconciliation); err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	if err := bird.ReconcileAll(ctx, r.Client, reconciliation); err != nil {
		return ctrl.Result{}, fmt.Errorf("error reconciling bird for routingv1alpha1.Router: %w", err)
	}

	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error storing ConfigMap: %w", err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RouterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&routingv1alpha1.Router{}).
		Owns(&corev1.ConfigMap{}).
		Watches(
			&routingv1alpha1.Session{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, a client.Object) []reconcile.Request {
				log := log.FromContext(ctx)

				session, ok := a.(*routingv1alpha1.Session)
				if !ok {
					log.WithValues("object", a).Info("Got unexepcted object while watching routingv1alpha1.Session")
					return nil
				}

				return []reconcile.Request{
					{NamespacedName: types.NamespacedName{
						Namespace: session.GetNamespace(),
						Name:      session.Spec.Router.Name,
					}},
				}
			}),
		).
		Complete(r)
}
