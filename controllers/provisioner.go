package controllers

import (
	"context"
	"github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/provisioner"
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
	ProvisionContainerImage = "quay.io/tmaxanc/hypersds-provisioner:v0.1.0"
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

	// TODO: SHOULD CHANGE THE LOGIC TO NOT USE PROVISIONER POD
	provisionorInstance, err := provisioner.NewProvisioner(r.Cluster.Spec, r.Client, r.Cluster.Namespace, r.Cluster.Name)
	if err != nil {
		return err
	}
	err = provisionorInstance.Run()
	if err != nil {
		return err
	}

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
