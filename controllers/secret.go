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

// SecretName indicates the name of the secret that contains admin keyring information
const SecretName = "ceph-secret"

func (r *CephClusterReconciler) syncSecret() error {
	if err := r.getSecret(); err == nil {
		return nil
	} else if !errors.IsNotFound(err) {
		return err
	}

	klog.Infof("syncSecret: creating secret %s", SecretName)
	updated, found, err := r.isAccessConfigMapUpdated()
	if err != nil {
		return err
	} else if !found {
		return nil
	}
	if updated {
		// This is the case when the secret is deleted after ceph cluster is completed.
		if err2 := r.updateAccessConfigMapAnnotation(false); err2 != nil {
			return err2
		}
	}

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
	return r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: SecretName}, secret)
}

func (r *CephClusterReconciler) newSecret() (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SecretName,
			Namespace: r.Cluster.Namespace,
		},
		Data: map[string][]byte{},
	}

	if err := controllerutil.SetControllerReference(r.Cluster, secret, r.Scheme); err != nil {
		return nil, err
	}
	return secret, nil
}
