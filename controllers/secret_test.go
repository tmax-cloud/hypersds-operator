package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("syncSecret", func() {
	Context("1. with no secret, updated=no", func() {
		cm := newConfigMap(AccessConfigMapName)
		cm.Annotations = map[string]string{
			"updated": "no",
		}
		r := createFakeCephClusterReconciler(cm)
		err := r.syncSecret()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should create a secret", func() {
			s := &corev1.Secret{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: SecretName}, s)
			Expect(err).Should(BeNil())
		})
	})

	Context("2. with no secret, updated=yes", func() {
		cm := newConfigMap(AccessConfigMapName)
		cm.Annotations = map[string]string{
			"updated": "yes",
		}
		r := createFakeCephClusterReconciler(cm)
		err := r.syncSecret()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should create a secret", func() {
			s := &corev1.Secret{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: SecretName}, s)
			Expect(err).Should(BeNil())
		})
		It("Should update configmap annotation(updated=no)", func() {
			cm := &corev1.ConfigMap{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: AccessConfigMapName}, cm)
			Expect(err).Should(BeNil())
			Expect(cm.Annotations["updated"]).Should(Equal("no"))
		})
	})

	Context("3. with secret", func() {
		s := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SecretName,
				Namespace: testCephClusterNs,
			},
		}
		r := createFakeCephClusterReconciler(s)
		err := r.syncSecret()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should not delete the secret", func() {
			s := &corev1.Secret{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: SecretName}, s)
			Expect(err).Should(BeNil())
		})
	})
})
