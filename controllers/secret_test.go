package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("syncSecret", func() {
	Context("1. without secret", func() {
		r := createFakeCephClusterReconciler()
		err := r.syncSecret()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should create a secret with suffix", func() {
			s := &corev1.Secret{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.getSecretName()}, s)
			Expect(err).Should(BeNil())
		})
		It("Should not create a secret without suffix", func() {
			s := &corev1.Secret{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.Cluster.Name}, s)
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("2. with secret", func() {
		s := newSecret()
		r := createFakeCephClusterReconciler(s)
		err := r.syncSecret()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should not delete the secret", func() {
			s := &corev1.Secret{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.getSecretName()}, s)
			Expect(err).Should(BeNil())
		})
	})
})
