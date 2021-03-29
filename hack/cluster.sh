#!/bin/bash

declare -A K8s_VMs
K8s_VMs=( ["pl-node1"]="192.168.56.31" \
         ["pl-node2"]="192.168.56.32" \
         ["pl-node3"]="192.168.56.33" )
declare -A CEPH_VMs
CEPH_VMs=( ["ubuntu-node1"]="192.168.33.11" \
             ["ubuntu-node2"]="192.168.33.12" \
             ["ubuntu-node3"]="192.168.33.13" )
VM="VirtualBox"

# Check if VirtualBox vm is running
is_vm_running(){
    VBoxManage list runningvms | grep "$1" &> /dev/null
}

# Check if VirtualBox vm exists
is_vm_exist(){
    VBoxManage list vms | grep "$1" &> /dev/null
}

# VirtualBox: Start vm
virtualbox_startvm(){
    echo "[$VM] start vm: '$1'..."
    if is_vm_exist "$1" && ! is_vm_running "$1"; then
        if ! VBoxManage startvm "$1" --type headless; then
            echo "[$VM] [ERROR] Unable to start '$1'!"
            exit 1
        fi
    else
        echo "[$VM] [ERROR] '$1' is NOT exist!"
        exit 1
    fi
}

# VirtualBox: stop vm
virtualbox_stopvm(){
    echo "[$VM] Poweroff vm: '$1'..."
    if is_vm_running "$1"; then
        if ! VBoxManage controlvm "$1" poweroff; then
            echo "[$VM] [ERROR] Unable to stop '$1'!"
            exit 1
        fi
    else
        echo "[$VM] $1 is not running!"
    fi
}

# VirtualBox: Check if snapshot exists
virtualbox_check_snapshot(){
    VBoxManage snapshot "$1" list | grep "$2" &> /dev/null
}

# VirtualBOx: restore snapshot
virtualbox_restore_snapshot(){
    echo "[$VM] '$1' restore snapshot: '$2'"
    if virtualbox_check_snapshot "$@"; then
        if is_vm_running "$1"; then
            virtualbox_stopvm "$1"
        fi

        if VBoxManage snapshot "$1" restore "$2"; then
            echo "[$VM] '$1' has been restored to: $2"
        else
            echo "[$VM] '$1' Unable to restore snapshot: $2"
            exit 1
        fi
    else
        echo "[$VM] [ERROR] '$1' snapshot ($2) is not exist!"
        exit 1
    fi
}

# Check if SSH to vm nodes is success
is_ssh_success(){
    local tries=1
    while ((tries < 30)); do
        echo "[SSH] Checking SSH connection to $1 : try ($tries)"
        if ssh -q root@"$1" exit; then
            echo "[SSH] SSH connection to $1 : SUCCESS"
            return 0
        fi
        tries=$((tries + 1))
        sleep 1
    done
    echo "[SSH] SSH connection to $1 : FAIL"
    return 1
}

is_all_pod_running(){
    local tries=1
    while ((tries < 30)); do
		echo "[Kubernetes] Waiting all pod to be Running..."
		if kubectl get node &> /dev/null; then
			if ! kubectl get pod -A | grep -E 'Error|CrashLoopBackOff' &> /dev/null; then
				kubectl cluster-info
				kubectl get node,pod -A -o wide
				return 0
			fi
		fi
        tries=$((tries + 1))
        sleep 3
    done
    return 1
}

# Copy Kubernetes admin.conf from master node to local node
get_kube_config(){
    KUBE_DIR=$1
    [ ! -d "$KUBE_DIR" ] && mkdir -p "$KUBE_DIR"
    for ip in "${K8s_VMs[@]}"; do
        if scp root@"${ip}":/etc/kubernetes/admin.conf "$KUBE_DIR"/config &> /dev/null; then
            sudo chown "$(id -u)":"$(id -g)" "$KUBE_DIR"/config &> /dev/null
            return 0
        fi
    done
    return 1
}

# Restart ntp service : solve ceph HEALTH_WARN (clock skew)
restart_ntp(){
    for ip in "${K8s_VMs[@]}"; do
		ssh root@"${ip}" 'systemctl restart chronyd' &> /dev/null
		ssh root@"${ip}" 'timedatectl set-ntp true' &> /dev/null
    done
}

# Loading requested cluster from snapshot
loading_kube(){
    for vm in "${!K8s_VMs[@]}"
    do
        virtualbox_restore_snapshot "$vm" "$1"
        echo
        virtualbox_startvm "$vm"
        echo
    done

    # Check SSH connection
    for ip in "${K8s_VMs[@]}"; do
        if ! is_ssh_success "$ip"; then
            exit 1
        fi
    done
}

# Loading requested cluster if exists
prolinuxKubeUp(){
    echo "[ClusterUP] Load Kubernetes cluster (Start)"
    kubeVer=$1
    kubeRun=$2
    kubeNet=$3

    if [ "$kubeVer" == "" ] ; then
        echo "[ERROR] Need env KUBE_VERSION=v1.xx.x"
        exit 1
    fi

    if [ "$kubeRun" == "" ] ; then
        echo "[INFO] Default container runtime (cri-o) is used!"
        kubeNet=crio
    fi

    if [ "$kubeNet" == "" ] ; then
        echo "[INFO] Default network plugin (calico) is used!"
        kubeNet=calico
    fi

    requestKube="${kubeVer}_${kubeRun}_${kubeNet}"
    echo "##### ($requestKube) is loading  #####"
    echo
    loading_kube "$requestKube"

    restart_ntp

    if get_kube_config "$HOME/.kube"; then
        export KUBECONFIG=$KUBE_DIR/config
		if ! is_all_pod_running; then
			echo "[ClusterUP] Cluster is not healthy."
			exit 1
		fi

    else
        echo "[ClusterUp] Unable to kube config 'admin.conf' from master node!"
        exit 1
    fi

    echo "[ClusterUP] Load Kubernetes cluster (Finish)"
    echo
}

# Clean kubernetes cluster: restore vm snapshot to No_k8s stage
prolinuxKubeClean(){
    loading_kube "no_k8s"
}

# Poweroff k8s nodes
prolinuxKubeDown(){
    for vm in "${!K8s_VMs[@]}"
    do
        virtualbox_stopvm "$vm"
    done
}

# Loading requested cluster from snapshot
loading_ceph(){
    for vm in "${!CEPH_VMs[@]}"
    do
        virtualbox_restore_snapshot "$vm" "$1"
        echo
        virtualbox_startvm "$vm"
        echo
    done

    # Check SSH connection
    for ip in "${CEPH_VMs[@]}"; do
        if ! is_ssh_success "$ip"; then
            exit 1
        fi
    done
}

cephNodeClean(){
    loading_ceph "initial_state"
}

# Poweroff ceph nodes
cephNodeDown(){
    for vm in "${!CEPH_VMs[@]}"
    do
        virtualbox_stopvm "$vm"
    done
}

# main
case "$1" in
    up)
        #KUBE_VERSION=v1.19.8 KUBE_RUNTIME=crio KUBE_NETWORK=calico
        prolinuxKubeUp "$KUBE_VERSION" "$KUBE_RUNTIME" "$KUBE_NETWORK"
	cephNodeClean
        ;;
    clean)
        prolinuxKubeClean
	cephNodeClean
        ;;
    down)
        prolinuxKubeDown
	cephNodeDown
        ;;
    *)
        echo "Usage: $0 COMMAND"
        echo "COMMAND:"
        echo "  up      : Load ProLinux kubernetes cluster with env KUBE_VERSION, KUBE_RUNTIME, KUBE_NETWORK"
        echo "  clean   : Clean ProLinux kubernetes cluster from nodes (no_k8s stage)"
        echo "  down    : Poweroff ProLinux kubernetes cluster nodes"
        echo 
        echo "Prerequisite:"
        echo "export KUBE_VERSION=v1.19.8 KUBE_RUNTIME=crio KUBE_NETWORK=calico"
        echo
        ;;
esac

