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
	if err := r.getConfigMap(); err == nil {
		return nil
	} else if !errors.IsNotFound(err) {
		return err
	}

	klog.Infof("syncConfigMap: creating config map %s", r.Cluster.Name)
	if err := r.updateState(v1alpha1.CephClusterStateCreating); err != nil {
		return err
	}
	if err := r.updateCondition(&metav1.Condition{
		Type:   string(v1alpha1.ConditionBootstrapped),
		Status: metav1.ConditionFalse,
		Reason: "BootstrappingIsNotFinished",
	}); err != nil {
		return err
	}
	if err := r.updateCondition(&metav1.Condition{
		Type:   string(v1alpha1.ConditionReadyToUse),
		Status: metav1.ConditionFalse,
		Reason: "CephClusterIsCreating",
	}); err != nil {
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

func (r *CephClusterReconciler) getConfigMap() error {
	cm := &corev1.ConfigMap{}
	if err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.getConfigMapName()}, cm); err != nil {
		return err
	}
	return nil
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
