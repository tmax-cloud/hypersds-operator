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

// RoleName indicates the name of the role for the provisioner pod
const RoleName = "provisioner-role"

func (r *CephClusterReconciler) syncRole() error {
	if err := r.getRole(); err == nil {
		return nil
	} else if !errors.IsNotFound(err) {
		return err
	}

	klog.Infof("syncRole: creating role %s", RoleName)
	newRole, err := r.newRole()
	if err != nil {
		return err
	}
	if err := r.Client.Create(context.TODO(), newRole); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (r *CephClusterReconciler) getRole() error {
	role := &corev1.Role{}
	return r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: RoleName}, role)
}

func (r *CephClusterReconciler) newRole() (*corev1.Role, error) {
	role := &corev1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RoleName,
			Namespace: r.Cluster.Namespace,
		},
		Rules: []corev1.PolicyRule{
			{
				Verbs:     []string{"get", "list", "update", "patch"},
				APIGroups: []string{""},
				Resources: []string{"configmaps", "secrets"},
			},
		},
	}

	if err := controllerutil.SetControllerReference(r.Cluster, role, r.Scheme); err != nil {
		return nil, err
	}
	return role, nil
}
