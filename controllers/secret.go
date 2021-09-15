package controllers

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const secretSuffix = "-keyring"

func (r *CephClusterReconciler) syncSecret() error {
	if err := r.getSecret(); err == nil {
		return nil
	} else if !errors.IsNotFound(err) {
		return err
	}

	klog.Infof("syncSecret: creating secret %s", r.Cluster.Name)
	newSecret, err := r.newSecret()
	if err != nil {
		return err
	}
	if err := r.Client.Create(context.TODO(), newSecret); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (r *CephClusterReconciler) getSecret() error {
	secret := &corev1.Secret{}
	if err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.getSecretName()}, secret); err != nil {
		return err
	}
	return nil
}

func (r *CephClusterReconciler) getSecretName() string {
	return r.Cluster.Name + secretSuffix
}

func (r *CephClusterReconciler) newSecret() (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.getSecretName(),
			Namespace: r.Cluster.Namespace,
		},
		Data: map[string][]byte{},
	}

	if err := controllerutil.SetControllerReference(r.Cluster, secret, r.Scheme); err != nil {
		return nil, err
	}
	return secret, nil
}
