package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("syncRoleBinding", func() {
	Context("1. with no role", func() {
		r := createFakeCephClusterReconciler()
		err := r.syncRoleBinding()

		It("Should return no error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should not create a rolebinding", func() {
			rb := &corev1.RoleBinding{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: RoleBindingName}, rb)
			Expect(errors.IsNotFound(err)).Should(BeTrue())
		})
	})

	Context("2. with role", func() {
		role := &corev1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      RoleName,
				Namespace: testCephClusterNs,
			},
		}
		r := createFakeCephClusterReconciler(role)
		err := r.syncRoleBinding()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should create a rolebinding", func() {
			rb := &corev1.RoleBinding{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: RoleBindingName}, rb)
			Expect(err).Should(BeNil())
		})
	})

	Context("3. with rolebinding", func() {
		rb := &corev1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      RoleBindingName,
				Namespace: testCephClusterNs,
			},
		}
		r := createFakeCephClusterReconciler(rb)
		err := r.syncRoleBinding()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should not delete the rolebinding", func() {
			rb := &corev1.RoleBinding{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: RoleBindingName}, rb)
			Expect(err).Should(BeNil())
		})
	})
})
