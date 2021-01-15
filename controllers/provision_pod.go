package controllers

import (
	"context"
	"github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	// ProvisionContainerName indicates container name of the provision pod
	ProvisionContainerName = "hypersds-provisioner"
	// ProvisionContainerImage indicates image name of the provision pod
	// TODO: ProvisionContainerImage should be updated
	ProvisionContainerImage = "192.168.7.16:5000/hypersds-provisioner"
	// ConfigMapVolumeName is used for creating configmap volume in pod specs
	ConfigMapVolumeName = "ceph-cluster-info"
	// ConfigMapVolumePath is a path where the configmap volume is mounted
	ConfigMapVolumePath = "/manifest"
)

func (r *CephClusterReconciler) syncProvisioner() error {
	if _, err := r.getConfigMap(getConfigMapName(r.Cluster.Name)); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	updated, foundCm, err := r.isAccessConfigMapUpdated()
	if err != nil {
		return err
	} else if !foundCm {
		return nil
	}

	if err = r.getSecret(); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	pod := &corev1.Pod{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: r.Cluster.Namespace, Name: getProvisionPodName(r.Cluster.Name)}, pod)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	foundPod := err == nil

	if !foundPod && !updated {
		klog.Infof("syncProvisioner: creating new pod for ceph cluster %s", r.Cluster.Name)
		newPod, err := r.newPod()
		if err != nil {
			return err
		}
		if err := r.Client.Create(context.TODO(), newPod); err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	} else if foundPod && isPodCompleted(pod) && !updated {
		klog.Infof("syncProvisioner: finishing for ceph cluster %s, deleting the pod and updating state", r.Cluster.Name)
		if err := r.updateAccessConfigMapAnnotation(true); err != nil {
			return err
		}
		if err := r.Client.Delete(context.TODO(), pod); err != nil && !errors.IsNotFound(err) {
			return err
		}
		if err := r.updateStateWithReadyToUse(v1alpha1.CephClusterStateCompleted, metav1.ConditionTrue, "CephClusterIsReady", "Ceph cluster is ready to use"); err != nil {
			return err
		}
	}
	return nil
}

func getProvisionPodName(cephClusterName string) string {
	return cephClusterName + "-provision-pod"
}

func isPodCompleted(pod *corev1.Pod) bool {
	return len(pod.Status.ContainerStatuses) != 0 &&
		pod.Status.ContainerStatuses[0].State.Terminated != nil &&
		pod.Status.ContainerStatuses[0].State.Terminated.Reason == "Completed"
}

func (r *CephClusterReconciler) newPod() (*corev1.Pod, error) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getProvisionPodName(r.Cluster.Name),
			Namespace: r.Cluster.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  ProvisionContainerName,
					Image: ProvisionContainerImage,
					Env: []corev1.EnvVar{
						{Name: "NAMESPACE", Value: r.Cluster.Namespace},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      ConfigMapVolumeName,
							MountPath: ConfigMapVolumePath,
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: ConfigMapVolumeName,
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: getConfigMapName(r.Cluster.Name),
							},
						},
					},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(r.Cluster, pod, r.Scheme); err != nil {
		return nil, err
	}
	return pod, nil
}
