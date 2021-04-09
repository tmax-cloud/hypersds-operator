package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("syncInstallConfig", func() {
	Context("1. with no configmap", func() {
		r := createFakeCephClusterReconciler()
		err := r.syncInstallConfig()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should create a configmap", func() {
			cm := &corev1.ConfigMap{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: getConfigMapName(r.Cluster.Name)}, cm)
			Expect(err).Should(BeNil())
		})
	})

	Context("2. with configmap", func() {
		cm := newConfigMap(getConfigMapName(testCephClusterName))
		r := createFakeCephClusterReconciler(cm)
		err := r.syncInstallConfig()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should not delete the configmap", func() {
			cm := &corev1.ConfigMap{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: getConfigMapName(r.Cluster.Name)}, cm)
			Expect(err).Should(BeNil())
		})
	})
})
