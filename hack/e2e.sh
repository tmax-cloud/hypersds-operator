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
    wait_condition "kubectl get cephclusters.hypersds.tmax.io | grep Completed" 1200 60
    kubectl delete -f config/samples/hypersds_v1alpha1_cephcluster.yaml
  ;;
  *)
    usage
  ;;
esac
