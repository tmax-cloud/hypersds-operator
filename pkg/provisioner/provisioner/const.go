package provisioner

// provisionerState is current installation phase of provisioner
type provisionerState string

const (
	// InstallStarted indicates Initialized phase
	InstallStarted provisionerState = "Initialized"
	// BasePkgInstalled indicates BaseInstalled phase
	BasePkgInstalled provisionerState = "BaseInstalled"
	// CephadmPkgInstalled indicates AdmInstalled phase
	CephadmPkgInstalled provisionerState = "AdmInstalled"
	// CephBootstrapped indicates Bootstrapped phase
	CephBootstrapped provisionerState = "Bootstrapped"
	// CephBootstrapCommitted indicates Committed phase
	CephBootstrapCommitted provisionerState = "Committed"
	// CephOsdDeployed indicates OsdDeployed phase
	CephOsdDeployed provisionerState = "OsdDeployed"
)

const (
	cephConfNameFromCr = "ceph_initial.conf"
	pathConfFromAdm    = "/etc/ceph/ceph.conf"
	pathKeyringFromAdm = "/etc/ceph/ceph.client.admin.keyring"
	cephVersion        = "15.2.8"
	cephImageName      = "ceph/ceph:v" + cephVersion
)
