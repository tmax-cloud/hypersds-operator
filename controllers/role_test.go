package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("syncRole", func() {
	Context("1. with no role", func() {
		r := createFakeCephClusterReconciler()
		err := r.syncRole()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should create a role", func() {
			role := &corev1.Role{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: RoleName}, role)
			Expect(err).Should(BeNil())
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
		err := r.syncRole()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should not delete the role", func() {
			role := &corev1.Role{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: RoleName}, role)
			Expect(err).Should(BeNil())
		})
	})
})
