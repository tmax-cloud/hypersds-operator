#!/bin/bash

function usage {
  echo " $0 [command]

  Available Commands:
    bootstrap
  " >&2
}

function wait_condition {
  cond=$1
  status=${cond%%|*}
  timeout=$2
  interval=$3

  for ((i=0; i<timeout; i+=$interval)) do
    if eval $cond > /dev/null 2>&1; then echo "Conditon met"; return 0; fi;
	    echo "Waiting $(((i+interval)/60)) minutes..."
    sleep $interval
    eval $status
  done

  echo "Timeout"
  return 1
}

case "${1:-}" in
  bootstrap)
    echo "deploying ceph cluster cr ..."
    kubectl apply -f config/samples/hypersds_v1alpha1_cephcluster.yaml
    wait_condition "kubectl get cephclusters.hypersds.tmax.io | grep Running" 2400 60
  ;;
  update_cm_after_delete)
    echo "deleting configmap ..."
    kubectl delete cm cephcluster-sample-conf
    sleep 3
    wait_condition "kubectl describe cm cephcluster-sample-conf | grep mon_host" 300 60
  ;;
  update_secret_after_delete)
    echo "deleting secret ..."
    kubectl delete secret cephcluster-sample-keyring
    sleep 3
    wait_condition "kubectl describe secret cephcluster-sample-keyring | grep keyring:" 300 60
  ;;
  add_host)
    echo "adding host ..."
    kubectl apply -f config/samples/hypersds_v1alpha1_cephcluster_hostadd.yaml
    sleep 600
  ;;
  add_disk)
    echo "adding disk ..."
    kubectl apply -f config/samples/hypersds_v1alpha1_cephcluster_diskadd_unordered.yaml
    sleep 600
  ;;
  delete_cluster)
    echo "deleting ceph cluster cr ..."
    kubectl delete -f config/samples/hypersds_v1alpha1_cephcluster.yaml
  ;;
  *)
    usage
  ;;
esac
