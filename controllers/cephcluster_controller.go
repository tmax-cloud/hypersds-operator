/*


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

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

// CephClusterReconciler reconciles a CephCluster object
type CephClusterReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	Cluster *hypersdsv1alpha1.CephCluster
}

// +kubebuilder:rbac:groups=hypersds.tmax.io,resources=cephclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hypersds.tmax.io,resources=cephclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps;secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile reads that state of the cluster for a CephCluster object and makes changes based on the state read
func (r *CephClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("cephcluster", req.NamespacedName)

	cachedCluster := &hypersdsv1alpha1.CephCluster{}
	if err := r.Client.Get(context.TODO(), req.NamespacedName, cachedCluster); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	r.Cluster = cachedCluster.DeepCopy()

	syncAll := func() error {
		if err := r.syncConfigMap(); err != nil {
			return err
		}
		if err := r.syncSecret(); err != nil {
			return err
		}
		if err := r.syncProvisioner(); err != nil {
			return err
		}
		return nil
	}
	if err := syncAll(); err != nil {
		if err2 := r.updateStateWithReadyToUse(hypersdsv1alpha1.CephClusterStateError, metav1.ConditionFalse, "SeeMessages", err.Error()); err2 != nil {
			return ctrl.Result{}, err2
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager creates a new CephCluster Controller and adds it to the Manager
func (r *CephClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hypersdsv1alpha1.CephCluster{}).
		Watches(&source.Kind{Type: &corev1.ConfigMap{}},
			&handler.EnqueueRequestForOwner{IsController: true, OwnerType: &hypersdsv1alpha1.CephCluster{}}).
		Watches(&source.Kind{Type: &corev1.Secret{}},
			&handler.EnqueueRequestForOwner{IsController: true, OwnerType: &hypersdsv1alpha1.CephCluster{}}).
		Complete(r)
}

func (r *CephClusterReconciler) updateStateWithReadyToUse(state hypersdsv1alpha1.CephClusterState, readyToUseStatus metav1.ConditionStatus,
	reason, message string) error {
	meta.SetStatusCondition(&r.Cluster.Status.Conditions, metav1.Condition{
		Type:    hypersdsv1alpha1.ConditionReadyToUse,
		Status:  readyToUseStatus,
		Reason:  reason,
		Message: message,
	})
	r.Cluster.Status.State = state
	return r.Client.Status().Update(context.TODO(), r.Cluster)
}
