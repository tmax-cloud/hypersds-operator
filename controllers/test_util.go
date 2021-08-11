package controllers

import (
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	testCephClusterName = "testcc"
	testCephClusterNs   = "default"
)

func createFakeCephClusterReconciler(objects ...runtime.Object) *CephClusterReconciler {
	cc := newCephCluster()
	c, s, err := createFakeClientAndScheme(append(objects, cc)...)
	if err != nil {
		panic(err)
	}
	return &CephClusterReconciler{Client: c, Scheme: s, Cluster: cc}
}

func newCephCluster() *hypersdsv1alpha1.CephCluster {
	return &hypersdsv1alpha1.CephCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testCephClusterName,
			Namespace: testCephClusterNs,
		},
		Spec: hypersdsv1alpha1.CephClusterSpec{
			Mon: hypersdsv1alpha1.CephClusterMonSpec{
				Count: 1,
			},
			Osd: []hypersdsv1alpha1.CephClusterOsdSpec{
				{
					HostName: "test",
					Devices:  []string{"/dev/sda", "/dev/sdb"},
				},
			},
			Nodes: []hypersdsv1alpha1.Node{
				{
					IP:       "0.0.0.0",
					UserID:   "test",
					Password: "test",
					HostName: "test",
				},
			},
		},
	}
}

func createFakeClientAndScheme(objects ...runtime.Object) (client.Client, *runtime.Scheme, error) {
	s := scheme.Scheme
	if err := hypersdsv1alpha1.AddToScheme(s); err != nil {
		return nil, nil, err
	}
	var objs []runtime.Object
	objs = append(objs, objects...)
	return fake.NewFakeClientWithScheme(s, objs...), s, nil
}

func newConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testCephClusterName + configMapSuffix,
			Namespace: testCephClusterNs,
		},
	}
}

func newSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testCephClusterName + secretSuffix,
			Namespace: testCephClusterNs,
		},
	}
}
