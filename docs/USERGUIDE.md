# User Guide

## Before you begin

- deploy K8s cluster through [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/), [kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/), [kubespray](https://kubernetes.io/docs/setup/production-environment/tools/kubespray/), or [other methods](https://kubernetes.io/docs/setup/)

> Note: Ceph CSI driver is needed to provision PVs and PVCs. A guide for how to deploy ceph csi driver and how to connect k8s cluster with external ceph cluster is located below this page.

## Table of Contents

- [Install Operator](#Install-Hypersds-Operator)
- [Bootstrap Ceph Cluster](#Use-Hypersds-Operator)
- [Uninstall Operator](#Uninstall-Hypersds-Operator)
- [Connect Ceph to K8s](#Use-external-ceph-cluster-in-Kubernetes)
- [Install Snapshot Controller to K8s](#Install-Snapshot-Controller)

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

## Use external ceph cluster in Kubernetes
This is a guide for connecting external ceph cluster to k8s using [Ceph CSI](https://github.com/ceph/ceph-csi). Ceph CSI plugins implement an interface between Kubernetes and Ceph cluster.
Independent CSI plugins are provided to support RBD and CephFS backed volumes. The Ceph CSI version used in this guide is [v3.1.2](https://github.com/ceph/ceph-csi/releases/tag/v3.1.2).

### 1. Gather the following information from the ceph cluster

#### Ceph cluster installed without HyperSDS-Operator (ex. cephadm or ansible)

##### 1-1. Change working directory on the node where the ceph cluster is installed
```shell script
$ cd /etc/ceph
```

##### 1-2. Ceph cluster fsid and monitor list
```shell script
$ cat ceph.conf
[global]
	fsid = 50611be6-33b3-11eb-a5cb-0894ef32cba4
	mon_host = v2:172.21.3.8:3300/0,v1:172.21.3.8:6789/0
```

##### 1-3. Key of the admin client
```shell script
$ cat ceph.client.admin.keyring
[client.admin]
	key = AQCeBcZftAEvExAAultsKBpNpiWWGi06Md7mmw==
```

##### 1-4. Name of the pool to be used for rbd
```shell script
$ ceph osd lspools
1 device_health_metrics
2 test-pool
3 rbd-pool
```

##### 1-5. Name of the volume to be used for cephfs
```shell script
$ ceph fs volume ls
[
    {
        "name": "myfs"
    }
]
```

#### Ceph cluster installed with HyperSDS-Operator
> Note: This part will be added soon.

### 2. Deploying CSI plugins with K8s
#### Deploy Namespace for CSI plugins and other resources
```shell script
$ kubectl apply -f deploy/csi/namespace.yaml
```

#### Deploy ConfigMap for CSI plugins
|Parameter |Value |
|---|---|
|`clusterID` | `fsid` from the step 1-2 |
|`monitors` | `mon_host` from the step 1-2 |

```shell script
# Replace `clusterID` and `monitors`.
$ kubectl apply -f deploy/csi/csi-config-map.yaml
```

#### Deploy Secret for CSI plugins
|Parameter |Value |
|---|---|
|`userKey` | `key` value from the step 1-3 |
|`adminKey` | `key` value from the step 1-3 |

```shell script
# rbd
# Replace `userKey`.
$ kubectl apply -f deploy/csi/rbd/secret.yaml

# cephfs
# Replace `adminKey`.
$ kubectl apply -f deploy/csi/cephfs/secret.yaml
```

#### Deploy RBAC
```shell script
# rbd
$ kubectl apply -f deploy/csi/rbd/csi-nodeplugin-rbac.yaml
$ kubectl apply -f deploy/csi/rbd/csi-provisioner-rbac.yaml

# cephfs
$ kubectl apply -f deploy/csi/cephfs/csi-nodeplugin-rbac.yaml
$ kubectl apply -f deploy/csi/cephfs/csi-provisioner-rbac.yaml
```

#### Deploy CSI plugins
```shell script
# rbd
$ kubectl apply -f deploy/csi/rbd/csi-rbdplugin-provisioner.yaml
$ kubectl apply -f deploy/csi/rbd/csi-rbdplugin.yaml

# verify deployment
$ kubectl get pod -n ceph-csi
NAMESPACE     NAME                                            READY   STATUS    RESTARTS   AGE
ceph-csi      csi-rbdplugin-fvh8j                             3/3     Running   0          28s
ceph-csi      csi-rbdplugin-provisioner-7646649999-4k8fz      6/6     Running   0          23s
ceph-csi      csi-rbdplugin-provisioner-7646649999-sc92b      6/6     Running   0          23s
ceph-csi      csi-rbdplugin-provisioner-7646649999-x2wg5      6/6     Running   0          23s

# cephfs
$ kubectl apply -f deploy/csi/cephfs/csi-cephfsplugin-provisioner.yaml
$ kubectl apply -f deploy/csi/cephfs/csi-cephfsplugin.yaml

# verify deployment
$ kubectl get pod -n ceph-csi
NAMESPACE     NAME                                            READY   STATUS    RESTARTS   AGE
ceph-csi      csi-cephfsplugin-k5mh5                          3/3     Running   0          51s
ceph-csi      csi-cephfsplugin-provisioner-66458c7db6-2sm8f   6/6     Running   0          43s
ceph-csi      csi-cephfsplugin-provisioner-66458c7db6-mcxms   6/6     Running   0          44s
ceph-csi      csi-cephfsplugin-provisioner-66458c7db6-swpqv   6/6     Running   0          43s
```

### 3. Verifying CSI plugins
#### Deploy StorageClass
|Parameter |Value |
|---|---|
|`clusterID` | `fsid` from the step 1-2 |
|`pool` | rbd pool name from the step 1-4 |
|`fsName` | cephfs volume name from the step 1-5 |

```shell script
# rbd
# Replace `clusterID` and `pool`.
$ kubectl apply -f deploy/csi/rbd/storageclass.yaml

$ kubectl get sc
NAME                 PROVISIONER                RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
csi-rbd-sc           rbd.csi.ceph.com           Delete          Immediate           true                   17s

# cephfs
# Replace `clusterID` and `fsName`.
$ kubectl apply -f deploy/csi/cephfs/storageclass.yaml

$ kubectl get sc
NAME                 PROVISIONER                RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
csi-cephfs-sc        cephfs.csi.ceph.com        Delete          Immediate           true                   2s
```

#### Deploy Pvc
```shell script
# rbd
$ kubectl apply -f deploy/csi/rbd/pvc.yaml

$ kubectl get pvc
NAME             STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
rbd-pvc          Bound    pvc-178a7e5b-1d64-485a-95bd-6980e9e0e793   1Gi        RWO            csi-rbd-sc      18m

# cephfs
$ kubectl apply -f deploy/csi/cephfs/pvc.yaml

$ kubectl get pvc
NAME             STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
csi-cephfs-pvc   Bound    pvc-02748c76-f985-4161-a191-b7ff4c030509   1Gi        RWX            csi-cephfs-sc   22m
```

#### Deploy Pod
```shell script
# rbd
$ kubectl apply -f deploy/csi/rbd/pod.yaml

$ kubectl get pod
NAME                  READY   STATUS    RESTARTS   AGE
busybox               1/1     Running   0          30s

# cephfs
$ kubectl apply -f deploy/csi/cephfs/pod.yaml

$ kubectl get pod
NAME                  READY   STATUS    RESTARTS   AGE
csi-cephfs-demo-pod   1/1     Running   0          25s
```

## Install Snapshot Controller
This is a guide for deploying snapshot controller using [external-snapshotter](https://github.com/kubernetes-csi/external-snapshotter/tree/v4.1.0).
Min K8s version is 1.20. The snapshot controller version used in this guide is v4.0.0.

### Before you begin
Install CSI provisioner through [Connect Ceph to K8s](#Use-external-ceph-cluster-in-Kubernetes).

### Deploy CRDs
```shell script
$ kubectl apply -f deploy/snapshot-controller/crd

# verify deployment
$ kubectl get crd
NAME                                                  CREATED AT
...
volumesnapshotclasses.snapshot.storage.k8s.io         2021-10-01T15:29:58Z
volumesnapshotcontents.snapshot.storage.k8s.io        2021-10-01T15:29:58Z
volumesnapshots.snapshot.storage.k8s.io               2021-10-01T15:29:58Z
```

### Deploy Snapshot Controller
```shell script
# rbac
$ kubectl apply -f deploy/snapshot-controller/rbac-snapshot-controller.yaml

# controller
$ kubectl apply -f deploy/snapshot-controller/setup-snapshot-controller.yaml

# verify deployment
$ kubectl get pod -A
NAMESPACE     NAME                                            READY   STATUS              RESTARTS   AGE
...
kube-system   snapshot-controller-9f68fdd9-k4lv9              1/1     Running             1          103m
kube-system   snapshot-controller-9f68fdd9-p4qx6              1/1     Running             0          103m
```

### Deploy VolumeSnapshotClass
- Replace `clusterID` with `fsid` in [ceph-conf](#1-2-ceph-cluster-fsid-and-monitor-list).
```shell script
# rbd
$ kubectl apply -f deploy/snapshot-controller/rbd/snapshotclass.yaml

# verify deployment
$ kubectl get volumesnapshotclass
NAME                   DRIVER             DELETIONPOLICY   AGE
csi-rbd-snapclass      rbd.csi.ceph.com   Delete           3s

# cephfs
$ kubectl apply -f deploy/snapshot-controller/cephfs/snapshotclass.yaml

# verify deployment
$ kubectl get volumesnapshotclass
NAME                      DRIVER                DELETIONPOLICY   AGE
csi-cephfs-snapclass      cephfs.csi.ceph.com   Delete           2s
```

### Create snapshot
- Replace `persistentVolumeClaimName` with the name of the pvc you want to create snapshot.
```shell script
# rbd
$ kubectl apply -f deploy/snapshot-controller/rbd/snapshot.yaml

# cephfs
$ kubectl apply -f deploy/snapshot-controller/cephfs/snapshot.yaml
```

### Create new pvc from snapshot
- Change the spec of pvc to the desired value.
```shell script
# rbd
$ kubectl apply -f deploy/snapshot-controller/rbd/restore.yaml

# cephfs
$ kubectl apply -f deploy/snapshot-controller/cephfs/restore.yaml
```
