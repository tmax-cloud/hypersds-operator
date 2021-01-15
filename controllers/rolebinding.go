package controllers

import (
	"context"
	corev1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// RoleBindingName indicates the name of the rolebinding for the provisioner pod
const RoleBindingName = "provisioner-rolebinding"

func (r *CephClusterReconciler) syncRoleBinding() error {
	if err := r.getRole(); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	rb := &corev1.RoleBinding{}
	if err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: RoleBindingName}, rb); err == nil {
		return nil
	} else if !errors.IsNotFound(err) {
		return err
	}

	klog.Infof("syncRoleBinding: creating rolebinding %s", RoleBindingName)
	newRb, err := r.newRoleBinding()
	if err != nil {
		return err
	}
	if err := r.Client.Create(context.TODO(), newRb); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (r *CephClusterReconciler) newRoleBinding() (*corev1.RoleBinding, error) {
	rb := &corev1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RoleBindingName,
			Namespace: r.Cluster.Namespace,
		},
		Subjects: []corev1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: r.Cluster.Namespace,
			},
		},
		RoleRef: corev1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     RoleName,
		},
	}

	if err := controllerutil.SetControllerReference(r.Cluster, rb, r.Scheme); err != nil {
		return nil, err
	}
	return rb, nil
}
