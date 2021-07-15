package controllers

import (
	"github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/provisioner"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

func (r *CephClusterReconciler) syncProvisioner() error {
	cmUpdated, err := r.isConfigMapUpdated()
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	secretUpdated, err := r.isSecretUpdated()
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if !cmUpdated || !secretUpdated {
		klog.Infof("syncProvisioner: bootstrapping ceph cluster %s", r.Cluster.Name)
		provisionerInstance, err := provisioner.NewProvisioner(r.Cluster.Spec, r.Client, r.Cluster.Namespace, r.Cluster.Name)
		if err != nil {
			return err
		}
		if err := provisionerInstance.Run(); err != nil {
			return err
		}
	}

	if err := r.updateStateWithReadyToUse(v1alpha1.CephClusterStateCompleted, v1.ConditionTrue, "CephClusterIsReady", "Ceph cluster is ready to use"); err != nil {
		return err
	}

	return nil
}
