# User Guide

## Before you begin

- deploy K8s cluster through [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/), [kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/), [kubespray](https://kubernetes.io/docs/setup/production-environment/tools/kubespray/), or [other methods](https://kubernetes.io/docs/setup/)

> Note: Ceph CSI driver is needed to provision PVs and PVCs. A guide for how to deploy ceph csi driver and how to connect k8s cluster with external ceph cluster is located below this page.

## Install Hypersds-Operator

> Note: hypersds-operator image will be available on this [public registry](https://quay.io/organization/tmaxanc) when the official release is ready. Until then local image registry is needed to deploy and use hypersds-operator.

``` shell
# Download hypersds-operator project 
$ git clone https://github.com/tmax-cloud/hypersds-operator.git
$ cd hypersds-operator/

# Deploy local docker image registry if you needed
$ make registry

# Local image build and push
$ make docker-build docker-push IMG=localhost:5000/hypersds-operator:v0.0.1

# Deploy CRD, operator, and etc.
$ make deploy IMG=localhost:5000/hypersds-operator:v0.0.1

# Check operator status
$ kubectl get deploy -n hypersds-operator-system
NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
hypersds-operator-controller-manager   1/1     1            1           11m
```

## Use Hypersds-Operator

> Note: Currently, hypersds-provisioner only bootstrap external ceph cluster. More functionalities will be available soon.

### Bootstrap ceph cluster

``` shell
# First, modify node information to bootstrap ceph cluster on the following yaml file
$ kubectl apply -f config/samples/hypersds_v1alpha1_cephcluster.yaml

# Wait until external ceph cluster is ready
$ kubectl get cephcluster
NAME       STATE
cephcluster-sample  Completed
```

## Uninstall Hypersds-Operator

> Note: Currently, hypersds-provisioner does not clean up external ceph cluster nodes and disks. Manual cleanup is needed to reset nodes. Uninstallation feature will be available soon.

``` shell
# Clean all deployed resources includes CRD
$ make clean

# Clean only CRD
$ make uninstall
```
