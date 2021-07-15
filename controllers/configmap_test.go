package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("syncConfigMap", func() {
	Context("1. with no configmap", func() {
		r := createFakeCephClusterReconciler()
		err := r.syncConfigMap()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should create a configmap with suffix", func() {
			cm := &corev1.ConfigMap{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.getConfigMapName()}, cm)
			Expect(err).Should(BeNil())
		})
		It("Should not create a configmap without suffix", func() {
			cm := &corev1.ConfigMap{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.Cluster.Name}, cm)
			Expect(err).ShouldNot(BeNil())
		})
		It("Should update state to Creating", func() {
			cc := &hypersdsv1alpha1.CephCluster{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.Cluster.Name}, cc)
			Expect(err).Should(BeNil())
			Expect(cc.Status.State).Should(Equal(hypersdsv1alpha1.CephClusterStateCreating))
		})
		It("Should update readyToUse to false", func() {
			cc := &hypersdsv1alpha1.CephCluster{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.Cluster.Name}, cc)
			Expect(err).Should(BeNil())
			cond := meta.FindStatusCondition(cc.Status.Conditions, hypersdsv1alpha1.ConditionReadyToUse)
			Expect(cond).ShouldNot(BeNil())
			Expect(cond.Status).Should(Equal(metav1.ConditionFalse))
		})
	})

	Context("2. with configmap", func() {
		cm := newConfigMap()
		r := createFakeCephClusterReconciler(cm)
		err := r.syncConfigMap()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should not delete the configmap", func() {
			cm := &corev1.ConfigMap{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.getConfigMapName()}, cm)
			Expect(err).Should(BeNil())
		})
	})
})

var _ = Describe("isConfigMapUpdated", func() {
	Context("1. with no configmap", func() {
		r := createFakeCephClusterReconciler()
		updated, err := r.isConfigMapUpdated()

		It("Should return error", func() {
			Expect(err).ShouldNot(BeNil())
		})
		It("Should return updated false", func() {
			Expect(updated).Should(BeFalse())
		})
	})

	Context("2. with empty configmap", func() {
		cm := newConfigMap()
		r := createFakeCephClusterReconciler(cm)
		updated, err := r.isConfigMapUpdated()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should return updated false", func() {
			Expect(updated).Should(BeFalse())
		})
	})

	Context("3. with configmap without ceph config key", func() {
		cm := newConfigMap()
		cm.Data = map[string]string{
			"test-key": "test-value",
		}
		r := createFakeCephClusterReconciler(cm)
		updated, err := r.isConfigMapUpdated()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should return updated false", func() {
			Expect(updated).Should(BeFalse())
		})
	})

	Context("4. with ceph config configmap", func() {
		cm := newConfigMap()
		cm.Data = map[string]string{
			"conf": "test-value",
		}
		r := createFakeCephClusterReconciler(cm)
		updated, err := r.isConfigMapUpdated()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should return updated true", func() {
			Expect(updated).Should(BeTrue())
		})
	})
})
