package controllers

import (
	"context"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *CephClusterReconciler) syncInstallConfig() error {
	if _, err := r.getConfigMap(getConfigMapName(r.Cluster.Name)); err == nil {
		return nil
	} else if !errors.IsNotFound(err) {
		return err
	}

	klog.Infof("syncInstallConfig: creating install config map %s", getConfigMapName(r.Cluster.Name))
	cm, err := r.newConfigMap()
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

func getConfigMapName(cephClusterName string) string {
	return cephClusterName + "-install-config"
}

func (r *CephClusterReconciler) newConfigMap() (*corev1.ConfigMap, error) {
	spec, err := yaml.Marshal(r.Cluster.Spec)
	if err != nil {
		return nil, err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getConfigMapName(r.Cluster.Name),
			Namespace: r.Cluster.Namespace,
		},
		Data: map[string]string{
			"cluster.yaml": string(spec),
		},
	}

	if err := controllerutil.SetControllerReference(r.Cluster, cm, r.Scheme); err != nil {
		return nil, err
	}
	return cm, nil
}
