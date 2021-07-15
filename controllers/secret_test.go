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

var _ = Describe("isSecretUpdated", func() {
	Context("1. with no secret", func() {
		r := createFakeCephClusterReconciler()
		updated, err := r.isSecretUpdated()

		It("Should return error", func() {
			Expect(err).ShouldNot(BeNil())
		})
		It("Should return updated false", func() {
			Expect(updated).Should(BeFalse())
		})
	})

	Context("2. with empty secret", func() {
		s := newSecret()
		r := createFakeCephClusterReconciler(s)
		updated, err := r.isSecretUpdated()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should return updated false", func() {
			Expect(updated).Should(BeFalse())
		})
	})

	Context("3. with secret without ceph keyring key", func() {
		s := newSecret()
		s.Data = map[string][]byte{
			"test-key": {65, 66, 67, 226, 130, 172},
		}
		r := createFakeCephClusterReconciler(s)
		updated, err := r.isSecretUpdated()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should return updated false", func() {
			Expect(updated).Should(BeFalse())
		})
	})

	Context("4. with ceph keyring secret", func() {
		s := newSecret()
		s.Data = map[string][]byte{
			"keyring": {65, 66, 67, 226, 130, 172},
		}
		r := createFakeCephClusterReconciler(s)
		updated, err := r.isSecretUpdated()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should return updated true", func() {
			Expect(updated).Should(BeTrue())
		})
	})
})
