package controllers

import (
	"context"
	"github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const configMapSuffix = "-conf"

func (r *CephClusterReconciler) syncConfigMap() error {
	// NotFound error will occur when the configmap is not created
	// No error will occur when the configmap is already created
	if _, err := r.getConfigMap(); err == nil {
		return nil
	} else if !errors.IsNotFound(err) {
		return err
	}

	klog.Infof("syncConfigMap: creating config map %s", r.Cluster.Name)
	if err := r.updateStateWithReadyToUse(v1alpha1.CephClusterStateCreating, metav1.ConditionFalse, "CephClusterIsCreating", "CephCluster is creating"); err != nil {
		return err
	}

	cm, err := r.newConfigMap()
	if err != nil {
		return err
	}
	if err := r.Client.Create(context.TODO(), cm); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (r *CephClusterReconciler) getConfigMap() (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	if err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.getConfigMapName()}, cm); err != nil {
		return nil, err
	}
	return cm, nil
}

func (r *CephClusterReconciler) getConfigMapName() string {
	return r.Cluster.Name + configMapSuffix
}

func (r *CephClusterReconciler) newConfigMap() (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.getConfigMapName(),
			Namespace: r.Cluster.Namespace,
		},
		Data: map[string]string{},
	}

	if err := controllerutil.SetControllerReference(r.Cluster, cm, r.Scheme); err != nil {
		return nil, err
	}
	return cm, nil
}

func (r *CephClusterReconciler) isConfigMapUpdated() (updated bool, err error) {
	cm, err := r.getConfigMap()
	if err != nil {
		return false, err
	}

	if cm.Data == nil {
		return false, nil
	}

	_, found := cm.Data["conf"]
	if !found {
		return false, nil
	}

	return true, nil
}
