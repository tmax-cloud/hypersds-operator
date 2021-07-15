package controllers

import (
	"context"
	goerrors "errors"
	"github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// AccessConfigMapName indicates the name of the configmap that contains ceph cluster access information
const AccessConfigMapName = "ceph-conf"

func (r *CephClusterReconciler) syncAccessConfig() error {
	// NotFound error will occur when the access configmap is not created
	// No error will occur when the access configmap is already created
	if _, err := r.getConfigMap(AccessConfigMapName); err == nil {
		return nil
	} else if !errors.IsNotFound(err) {
		return err
	}

	klog.Infof("syncAccessConfig: creating access config map %s", AccessConfigMapName)
	if err := r.updateStateWithReadyToUse(v1alpha1.CephClusterStateCreating, metav1.ConditionFalse, "CephClusterIsCreating", "CephCluster is creating"); err != nil {
		return err
	}

	cm, err := r.newAccessConfigMap()
	if err != nil {
		return err
	}
	if err := r.Client.Create(context.TODO(), cm); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (r *CephClusterReconciler) getConfigMap(name string) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	if err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: name}, cm); err != nil {
		return nil, err
	}
	return cm, nil
}

func (r *CephClusterReconciler) newAccessConfigMap() (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        AccessConfigMapName,
			Namespace:   r.Cluster.Namespace,
			Annotations: map[string]string{"updated": "no"},
		},
		Data: map[string]string{},
	}

	if err := controllerutil.SetControllerReference(r.Cluster, cm, r.Scheme); err != nil {
		return nil, err
	}
	return cm, nil
}

func (r *CephClusterReconciler) isAccessConfigMapUpdated() (updated, found bool, err error) {
	cm, err := r.getConfigMap(AccessConfigMapName)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, false, nil
		}
		return false, false, err
	}

	u, found := cm.Annotations["updated"]
	if !found {
		return false, true, goerrors.New("invalid configmap annotation: Must need 'updated'")
	}
	return u == "yes", true, nil
}

func (r *CephClusterReconciler) updateAccessConfigMapAnnotation(updated bool) error {
	cm, err := r.getConfigMap(AccessConfigMapName)
	if err != nil {
		return err
	}
	if cm.Annotations == nil {
		cm.Annotations = map[string]string{}
	}
	if updated {
		cm.Annotations["updated"] = "yes"
	} else {
		cm.Annotations["updated"] = "no"
	}
	return r.Client.Update(context.TODO(), cm)
}
