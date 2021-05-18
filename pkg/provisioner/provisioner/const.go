package provisioner

// ProvisionerState is current installation phase of provisioner
type ProvisionerState string

const (
	InstallStarted         ProvisionerState = "Initialized"
	BasePkgInstalled       ProvisionerState = "BaseInstalled"
	CephadmPkgInstalled    ProvisionerState = "AdmInstalled"
	CephBootstrapped       ProvisionerState = "Bootstrapped"
	CephBootstrapCommitted ProvisionerState = "Committed"
	CephOsdDeployed        ProvisionerState = "OsdDeployed"
)

// File name, path, etc
const (
	cephConfNameFromCr = "ceph_initial.conf"
	pathConfFromAdm    = "/etc/ceph/ceph.conf"
	pathKeyringFromAdm = "/etc/ceph/ceph.client.admin.keyring"
	cephVersion        = "15.2.8"
	cephImageName      = "ceph/ceph:v" + cephVersion
)
