package util

const (
	// PathGlobalConfigDir is global directory path of ceph config
	PathGlobalConfigDir = "/working/config/"
	// CephConfName is filename of ceph conf file
	CephConfName = "ceph.conf"
	// CephKeyringName is filename of ceph keyring file
	CephKeyringName = "ceph.client.admin.keyring"
	// K8sConfigMapSuffix is suffix for the k8s configmap name
	K8sConfigMapSuffix = "-conf"
	// K8sSecretSuffix is suffix for the k8s secret name
	K8sSecretSuffix = "-keyring"
)
