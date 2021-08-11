package controllers

//	no.	configmap	secret	cm-updated	secret-updated	ready-to-use
//	1	x	x
//	2	o	x
//	3	x	o
//	4	o	o		o		o	true

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("syncProvisioner", func() {
	Context("1. without any resources", func() {
		r := createFakeCephClusterReconciler()
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("2. without secret", func() {
		cm := newConfigMap()
		r := createFakeCephClusterReconciler(cm)
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("3. without configmap", func() {
		s := newSecret()
		r := createFakeCephClusterReconciler(s)
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("4. with updated configmap and updated secret", func() {
		cm := newConfigMap()
		cm.Data = map[string]string{
			"conf": "no",
		}
		s := newSecret()
		s.Data = map[string][]byte{
			"keyring": {65, 66, 67, 226, 130, 172},
		}

		r := createFakeCephClusterReconciler(cm, s)
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should update state to Completed", func() {
			cc := &v1alpha1.CephCluster{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.Cluster.Name}, cc)
			Expect(err).Should(BeNil())
			Expect(cc.Status.State).Should(Equal(v1alpha1.CephClusterStateCompleted))
		})
		It("Should update readyToUse to true", func() {
			cc := &v1alpha1.CephCluster{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.Cluster.Name}, cc)
			Expect(err).Should(BeNil())
			Expect(meta.IsStatusConditionTrue(cc.Status.Conditions, v1alpha1.ConditionReadyToUse)).Should(BeTrue())
		})
	})
})
