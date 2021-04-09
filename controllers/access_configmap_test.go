package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("syncAccessConfig", func() {
	Context("1. with no configmap", func() {
		r := createFakeCephClusterReconciler()
		err := r.syncAccessConfig()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should create a configmap", func() {
			cm := &corev1.ConfigMap{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: AccessConfigMapName}, cm)
			Expect(err).Should(BeNil())
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
		cm := newConfigMap(AccessConfigMapName)
		r := createFakeCephClusterReconciler(cm)
		err := r.syncAccessConfig()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should not delete the configmap", func() {
			cm := &corev1.ConfigMap{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: AccessConfigMapName}, cm)
			Expect(err).Should(BeNil())
		})
	})
})

var _ = Describe("isAccessConfigMapUpdated", func() {
	Context("1. with no configmap", func() {
		r := createFakeCephClusterReconciler()
		_, found, err := r.isAccessConfigMapUpdated()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should return not found", func() {
			Expect(found).Should(BeFalse())
		})
	})

	Context("2. with configmap without annotation", func() {
		cm := newConfigMap(AccessConfigMapName)
		cm.Annotations = nil
		r := createFakeCephClusterReconciler(cm)
		_, _, err := r.isAccessConfigMapUpdated()

		It("Should return error", func() {
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("3. with configmap with annotation(updated=no)", func() {
		cm := newConfigMap(AccessConfigMapName)
		cm.Annotations = map[string]string{
			"updated": "no",
		}
		r := createFakeCephClusterReconciler(cm)
		u, found, err := r.isAccessConfigMapUpdated()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should return updated false", func() {
			Expect(u).Should(BeFalse())
		})
		It("Should return found configmap", func() {
			Expect(found).Should(BeTrue())
		})
	})

	Context("4. with configmap with annotation(updated=yes)", func() {
		cm := newConfigMap(AccessConfigMapName)
		cm.Annotations = map[string]string{
			"updated": "yes",
		}
		r := createFakeCephClusterReconciler(cm)
		u, found, err := r.isAccessConfigMapUpdated()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should return updated true", func() {
			Expect(u).Should(BeTrue())
		})
		It("Should return found configmap", func() {
			Expect(found).Should(BeTrue())
		})
	})
})

var _ = Describe("updateAccessConfigMapAnnotation", func() {
	Context("1. with no configmap", func() {
		r := createFakeCephClusterReconciler()
		err := r.updateAccessConfigMapAnnotation(true)

		It("Should return not found error", func() {
			Expect(errors.IsNotFound(err)).Should(BeTrue())
		})
	})

	Context("2. with configmap and updated to no", func() {
		cm := newConfigMap(AccessConfigMapName)
		r := createFakeCephClusterReconciler(cm)
		err := r.updateAccessConfigMapAnnotation(false)

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should configmap has annotation with updated no", func() {
			cm := &corev1.ConfigMap{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: AccessConfigMapName}, cm)
			Expect(err).Should(BeNil())
			Expect(cm.Annotations["updated"]).Should(Equal("no"))
		})
	})

	Context("3. with configmap and updated to yes", func() {
		cm := newConfigMap(AccessConfigMapName)
		r := createFakeCephClusterReconciler(cm)
		err := r.updateAccessConfigMapAnnotation(true)

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should configmap has annotation with updated yes", func() {
			cm := &corev1.ConfigMap{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: AccessConfigMapName}, cm)
			Expect(err).Should(BeNil())
			Expect(cm.Annotations["updated"]).Should(Equal("yes"))
		})
	})
})
