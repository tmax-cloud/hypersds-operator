package controllers

import (
	"context"
	"github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/provisioner"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

func (r *CephClusterReconciler) syncProvisioner() error {
	if err := r.getConfigMap(); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if err := r.getSecret(); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if err := r.updateDeployNode(); err != nil {
		return err
	}

	klog.Infof("syncProvisioner: bootstrapping ceph cluster %s", r.Cluster.Name)
	provisionerInstance, err := provisioner.NewProvisioner(r.Cluster.Spec, r.Client, r.Cluster.Namespace, r.Cluster.Name)
	if err != nil {
		return err
	}
	if err := provisionerInstance.Run(); err != nil {
		return err
	}
	if err := r.updateState(v1alpha1.CephClusterStateRunning); err != nil {
		return err
	}
	if err := r.updateCondition(&v1.Condition{
		Type:   string(v1alpha1.ConditionReadyToUse),
		Status: v1.ConditionTrue,
		Reason: "CephClusterIsReady",
	}); err != nil {
		return err
	}
	if meta.IsStatusConditionFalse(r.Cluster.Status.Conditions, string(v1alpha1.ConditionBootstrapped)) {
		if err := r.updateCondition(&v1.Condition{
			Type:   string(v1alpha1.ConditionBootstrapped),
			Status: v1.ConditionTrue,
			Reason: "CephClusterIsBootstrapped",
		}); err != nil {
			return err
		}
	}

	return nil
}

func (r *CephClusterReconciler) updateDeployNode() error {
	if r.Cluster.Status.DeployNode == (v1alpha1.Node{}) {
		r.Cluster.Status.DeployNode = r.getDefaultDeployNode()
		if err := r.Client.Status().Update(context.TODO(), r.Cluster); err != nil {
			return err
		}
	}
	return nil
}

// Decide deploying node (currently, first node is deploying node)
func (r *CephClusterReconciler) getDefaultDeployNode() v1alpha1.Node {
	return r.Cluster.Spec.Nodes[0]
}
