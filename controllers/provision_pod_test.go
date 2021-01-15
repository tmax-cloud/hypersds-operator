package controllers

//	no.	config	accessConfig	secret	pod	updated	state
//	1	x	x		x	x
//	2	o	x		x	x
//	3	o	o		x	x
//	4	o	o		o	x	no
//	5	o	o		o	x	yes
//	6	o	o		o	o	no	Completed
//	7	o	o		o	o	no	Running

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

var _ = Describe("syncProvisioner", func() {
	Context("1. with no resource", func() {
		r := createFakeCephClusterReconciler()
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("2. with configmap", func() {
		cm := newConfigMap(getConfigMapName(testCephClusterName))
		r := createFakeCephClusterReconciler(cm)
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("3. with configmap and access configmap", func() {
		cm := newConfigMap(getConfigMapName(testCephClusterName))
		accessCm := newConfigMap(AccessConfigMapName)
		accessCm.Annotations = map[string]string{
			"updated": "no",
		}
		r := createFakeCephClusterReconciler(cm, accessCm)
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("4. with configmap, access configmap and secret, updated=no", func() {
		cm := newConfigMap(getConfigMapName(testCephClusterName))
		accessCm := newConfigMap(AccessConfigMapName)
		accessCm.Annotations = map[string]string{
			"updated": "no",
		}
		s := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SecretName,
				Namespace: testCephClusterNs,
			},
		}
		r := createFakeCephClusterReconciler(cm, accessCm, s)
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should create a pod", func() {
			pod := &corev1.Pod{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: getProvisionPodName(r.Cluster.Name)}, pod)
			Expect(err).Should(BeNil())
		})
	})

	Context("5. with configmap, access configmap and secret, updated=yes", func() {
		cm := newConfigMap(getConfigMapName(testCephClusterName))
		accessCm := newConfigMap(AccessConfigMapName)
		accessCm.Annotations = map[string]string{
			"updated": "yes",
		}
		s := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SecretName,
				Namespace: testCephClusterNs,
			},
		}
		r := createFakeCephClusterReconciler(cm, accessCm, s)
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("6. with configmap, access configmap, secret and pod(Completed), updated=no", func() {
		cm := newConfigMap(getConfigMapName(testCephClusterName))
		accessCm := newConfigMap(AccessConfigMapName)
		accessCm.Annotations = map[string]string{
			"updated": "no",
		}
		s := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SecretName,
				Namespace: testCephClusterNs,
			},
		}
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      getProvisionPodName(testCephClusterName),
				Namespace: testCephClusterNs,
			},
			Status: corev1.PodStatus{
				ContainerStatuses: []corev1.ContainerStatus{
					{
						State: corev1.ContainerState{
							Terminated: &corev1.ContainerStateTerminated{
								Reason: "Completed",
							},
						},
					},
				},
			},
		}
		r := createFakeCephClusterReconciler(cm, accessCm, s, pod)
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should configmap has annotation with updated yes", func() {
			cm := &corev1.ConfigMap{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: AccessConfigMapName}, cm)
			Expect(err).Should(BeNil())
			Expect(cm.Annotations["updated"]).Should(Equal("yes"))
		})
		It("Should delete the pod", func() {
			pod := &corev1.Pod{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: getProvisionPodName(r.Cluster.Name)}, pod)
			Expect(errors.IsNotFound(err)).Should(BeTrue())
		})
		It("Should update state to Completed", func() {
			cc := &hypersdsv1alpha1.CephCluster{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.Cluster.Name}, cc)
			Expect(err).Should(BeNil())
			Expect(cc.Status.State).Should(Equal(hypersdsv1alpha1.CephClusterStateCompleted))
		})
		It("Should update readyToUse to true", func() {
			cc := &hypersdsv1alpha1.CephCluster{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: r.Cluster.Name}, cc)
			Expect(err).Should(BeNil())
			cond := meta.FindStatusCondition(cc.Status.Conditions, hypersdsv1alpha1.ConditionReadyToUse)
			Expect(cond).ShouldNot(BeNil())
			Expect(cond.Status).Should(Equal(metav1.ConditionTrue))
		})
	})

	Context("7. with configmap, access configmap, secret and pod(Running), updated=no", func() {
		cm := newConfigMap(getConfigMapName(testCephClusterName))
		accessCm := newConfigMap(AccessConfigMapName)
		accessCm.Annotations = map[string]string{
			"updated": "no",
		}
		s := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SecretName,
				Namespace: testCephClusterNs,
			},
		}
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      getProvisionPodName(testCephClusterName),
				Namespace: testCephClusterNs,
			},
		}
		r := createFakeCephClusterReconciler(cm, accessCm, s, pod)
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
		It("Should not delete the pod", func() {
			pod := &corev1.Pod{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: getProvisionPodName(r.Cluster.Name)}, pod)
			Expect(err).Should(BeNil())
		})
	})
})
