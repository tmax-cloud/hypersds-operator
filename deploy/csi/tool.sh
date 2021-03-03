#!/bin/bash

function config_install() {
  kubectl apply -f namespace.yaml
  kubectl apply -f csi-config-map.yaml
}

function cephfs_install() {
  kubectl apply -f cephfs/csi-nodeplugin-rbac.yaml
  kubectl apply -f cephfs/csi-provisioner-rbac.yaml
  kubectl apply -f cephfs/csi-cephfsplugin-provisioner.yaml
  kubectl apply -f cephfs/csi-cephfsplugin.yaml
  kubectl apply -f cephfs/secret.yaml
}

function cephfs_test() {
  kubectl apply -f cephfs/storageclass.yaml
  kubectl apply -f cephfs/pvc.yaml
  kubectl apply -f cephfs/pod.yaml
}

function rbd_install() {
  kubectl apply -f rbd/csi-nodeplugin-rbac.yaml
  kubectl apply -f rbd/csi-provisioner-rbac.yaml
  kubectl apply -f rbd/csi-rbdplugin-provisioner.yaml
  kubectl apply -f rbd/csi-rbdplugin.yaml
  kubectl apply -f rbd/secret.yaml
}

function rbd_test() {
  kubectl apply -f rbd/storageclass.yaml
  kubectl apply -f rbd/pvc.yaml
  kubectl apply -f rbd/pod.yaml
}

function clean() {
  kubectl delete -f cephfs/pod.yaml
  kubectl delete -f cephfs/pvc.yaml
  kubectl delete -f rbd/pod.yaml
  kubectl delete -f rbd/pvc.yaml
  kubectl delete -f namespace.yaml
  kubectl delete -f cephfs/storageclass.yaml
  kubectl delete -f rbd/storageclass.yaml
}

case "${1:-}" in
config_install)
  config_install
  ;;
cephfs_install)
  cephfs_install
  ;;
rbd_install)
  rbd_install
  ;;
cephfs_test)
  cephfs_test
  ;;
rbd_test)
  rbd_test
  ;;
clean)
  clean
  ;;
*)
  print_help
  ;;
esac
