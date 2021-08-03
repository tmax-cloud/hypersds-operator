package provisioner

import (
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/util"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/config"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/node"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"context"
	"errors"
	"fmt"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Provisioner contains information for ceph deployment and current ceph deployment status
type Provisioner struct {
	cephCluster hypersdsv1alpha1.CephClusterSpec
	state       provisionerState

	cephNamespace string
	cephName      string
	clientSet     client.Client
	cephConfig    *config.CephConfig
}

func (p *Provisioner) getState() provisionerState {
	return p.state
}

func (p *Provisioner) getNodes() ([]*node.Node, error) {
	return node.NewNodesFromCephCr(p.cephCluster)
}
func (p *Provisioner) getCephName() string {
	return p.cephName
}
func (p *Provisioner) getCephNamespace() string {
	return p.cephNamespace
}
func (p *Provisioner) getCephConfig() *config.CephConfig {
	return p.cephConfig
}
func (p *Provisioner) getClientSet() client.Client {
	return p.clientSet
}
func (p *Provisioner) getPathConfigDir() string {
	return util.PathGlobalConfigDir + p.cephName + "/"
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (p *Provisioner) setState(state provisionerState) error {
	p.state = state
	return nil
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (p *Provisioner) setCephCluster(cephCluster hypersdsv1alpha1.CephClusterSpec) error {
	p.cephCluster = cephCluster
	return nil
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (p *Provisioner) setCephNamespace(cephNamespace string) error {
	p.cephNamespace = cephNamespace
	return nil
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (p *Provisioner) setCephName(cephName string) error {
	p.cephName = cephName
	return nil
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (p *Provisioner) setClientSet(clientSet client.Client) error {
	p.clientSet = clientSet
	return nil
}

func (p *Provisioner) setCephConfig() error {
	var err error
	p.cephConfig, err = config.NewConfigFromCephCr(p.cephCluster)
	return err
}

// Run executes the following ceph deployment steps according to the current ceph deployment state
func (p *Provisioner) Run() error {
	// Decide deploying node (currently, first node is deploying node)
	nodeList, err := p.getNodes()
	if err != nil {
		return err
	}
	deployNode := nodeList[0]

	// Create config object from Ceph CR
	err = p.setCephConfig()
	if err != nil {
		return err
	}
	cephConfig := p.getCephConfig()
	pathConfigDir := p.getPathConfigDir()
	pathConfFromCr := pathConfigDir + cephConfNameFromCr
	pathConf := pathConfigDir + util.CephConfName
	pathKeyring := pathConfigDir + util.CephKeyringName

	switch currentState := p.getState(); currentState {
	case InstallStarted:
		// Fetch OS information for each nodes
		err = fetchOSInfo(nodeList)
		if err != nil {
			return err
		}

		// Install base package to all nodes
		err = installBasePackage(nodeList)
		if err != nil {
			return err
		}

		// Set provisioner state to BasePkgInstalled
		err = p.setState(BasePkgInstalled)
		if err != nil {
			return err
		}

		fallthrough

	case BasePkgInstalled:
		// Install cephadm package to deploying node
		err = installCephadm(deployNode)
		if err != nil {
			return err
		}

		// Set provisioner state to CephadmPkgInstalled
		err = p.setState(CephadmPkgInstalled)
		if err != nil {
			return err
		}

		fallthrough

	case CephadmPkgInstalled:
		// Extract initial conf file of Ceph
		err = cephConfig.MakeIniFile(wrapper.IoUtilWrapper, pathConfFromCr)
		if err != nil {
			return err
		}

		// Bootstrap ceph on deploy node with cephadm
		err = bootstrapCephadm(deployNode, pathConfFromCr)
		if err != nil {
			return err
		}

		// Set provisioner state to CephBootstrapped
		err = p.setState(CephBootstrapped)
		if err != nil {
			return err
		}

		fallthrough

	case CephBootstrapped:
		// Copy conf and keyring from deploy node
		err = copyFile(deployNode, node.SOURCE, pathConfFromAdm, pathConf)
		if err != nil {
			return err
		}
		err = copyFile(deployNode, node.SOURCE, pathKeyringFromAdm, pathKeyring)
		if err != nil {
			return err
		}

		// Update conf and keyring to ConfigMap and Secret
		err = p.updateCephClusterToOp()
		if err != nil {
			return err
		}

		// Set provisioner state to CephBootstrapCommitted
		err = p.setState(CephBootstrapCommitted)
		if err != nil {
			return err
		}

		fallthrough

	case CephBootstrapCommitted:
		// Copy conf and keyring from deploy node for ceph-common
		cephConf, cephKeyring, _, _ := p.checkKubeObjectUpdated()
		// todo check logic?

		err = p.applyHost(wrapper.YamlWrapper, wrapper.ExecWrapper, wrapper.IoUtilWrapper, cephConf, cephKeyring)
		if err != nil {
			return err
		}

		err = p.applyOsd(cephConf, cephKeyring)
		if err != nil {
			return err
		}
		err = p.setState(CephOsdDeployed)
		if err != nil {
			return err
		}
	}
	err = wrapper.OsWrapper.RemoveAll(pathConfigDir)
	if err != nil {
		fmt.Println("[Provisioner] CephName Driectory Remove Error")
		return err
	}

	return nil
}

func (p *Provisioner) identifyProvisionerState() (provisionerState, error) {
	// Decide deploying node (currently, first node is deploying node)
	nodes, err := p.getNodes()
	if err != nil {
		return "", err
	}
	deployNode := nodes[0]

	// Check base pkgs are installed
	// TODO: May contain error if user removed docker but did not purge dpkg
	const checkDockerWorkingCmd = "docker -v"
	_, err = deployNode.RunSSHCmd(wrapper.SSHWrapper, checkDockerWorkingCmd)
	// It considers any error that base pkgs are not installed
	if err != nil {
		// TODO: Replace stdout to log out
		fmt.Println("[identifyProvisionerState] docker is not installed")
		return InstallStarted, nil
	}

	// Check Cephadm is installed
	const checkCephadmInstalledCmd = "cephadm version"
	outputCephadm, err := deployNode.RunSSHCmd(wrapper.SSHWrapper, checkCephadmInstalledCmd)
	if err != nil {
		const cephadmNotFound = "command not found"
		if strings.Contains(outputCephadm.String(), cephadmNotFound) {
			return BasePkgInstalled, nil
		}

		// Other error occurred on RunSSHCmd
		// TODO: Replace stdout to log out
		fmt.Println("[identifyProvisionerState] cephadm installation check is failed")
		return BasePkgInstalled, err
	}

	// Check Ceph is bootstrapped
	const checkBootstrappedCmd = "cephadm shell -- ceph -s"
	outputBootstrap, err := deployNode.RunSSHCmd(wrapper.SSHWrapper, checkBootstrappedCmd)
	if err != nil {
		const objectNotFound = "ObjectNotFound"
		if strings.Contains(outputBootstrap.String(), objectNotFound) {
			return CephadmPkgInstalled, nil
		}

		// Other error occurred on RunSSHCmd
		// TODO: Replace stdout to log out
		fmt.Println("[identifyProvisionerState] ceph bootstrap check is failed")
		return CephadmPkgInstalled, err
	}

	// Confirm some Ceph health is returned
	const cephHealthPrefix = "HEALTH_"
	if !strings.Contains(outputBootstrap.String(), cephHealthPrefix) {
		// TODO: Replace stdout to log out
		fmt.Println("[identifyProvisionerState] ceph status does not return HEALTH_*")

		// TODO: Make own error pkg of hypersds-provisioner
		return CephadmPkgInstalled,
			errors.New("Error on Ceph bootstrap, cech status result: \n" +
				outputBootstrap.String())
	}

	// Check Ceph bootstrap is committed
	_, _, committed, err := p.checkKubeObjectUpdated()
	if err != nil {
		// Other error occurred on checkKubeObjectUpdated
		// TODO: Replace stdout to log out
		fmt.Println("[identifyProvisionerState] k8s configmap and secret check is failed")
		return CephBootstrapped, err
	}

	if !committed {
		// TODO: Replace stdout to log out
		fmt.Println("[identifyProvisionerState] k8s configmap and secret are not updated")
		return CephBootstrapped, nil
	}

	return CephBootstrapCommitted, nil
}

// TODO: Replace config const to inputs (e.g. K8sConfigMap, etc)
func (p *Provisioner) checkKubeObjectUpdated() (confDataBuf, keyringDataBuf []byte, phase bool, err error) {
	// Check ceph.conf is updated to ConfigMap
	configMap := &corev1.ConfigMap{}
	if err := p.clientSet.Get(context.TODO(), types.NamespacedName{Namespace: p.cephNamespace, Name: p.cephName + util.K8sConfigMapSuffix}, configMap); err != nil {
		// Configmap must exist
		if kubeerrors.IsNotFound(err) {
			// TODO: Replace stdout to log out
			fmt.Println("ConfigMap must exist")
			return nil, nil, false, err
		}
		return nil, nil, false, err
	}

	// Bootstrap commit has not occurred
	if configMap.Data == nil {
		return nil, nil, false, nil
	}
	confData, ok := configMap.Data["conf"]
	if !ok {
		fmt.Println("Ceph Conf Data isn't exist")
		return nil, nil, false, errors.New("Ceph Conf Data isn't exist")
	}
	confDataBuf = []byte(confData)

	// Check client.admin.keyring is updated to Secret
	secret := &corev1.Secret{}
	if err := p.clientSet.Get(context.TODO(), types.NamespacedName{Namespace: p.cephNamespace, Name: p.cephName + util.K8sSecretSuffix}, secret); err != nil {
		if kubeerrors.IsNotFound(err) {
			// TODO: Replace stdout to log out
			fmt.Println("Secret must exist")
			return nil, nil, false, err
		}
		return nil, nil, false, err
	}

	if secret.Data == nil {
		return nil, nil, false, nil
	}

	keyringDataBuf, ok = secret.Data["keyring"]
	if !ok {
		fmt.Println("Ceph Keyring Data isn't exist")
		return nil, nil, false, errors.New("Ceph keyring Data isn't exist")
	}

	return confDataBuf, keyringDataBuf, true, nil
}

// NewProvisioner creates Provisioner using ceph deployment information and checks the current ceph deployment status in node
func NewProvisioner(cephClusterSpec hypersdsv1alpha1.CephClusterSpec, clientSet client.Client, cephNamespace, cephName string) (*Provisioner, error) {
	var err error
	var currentState provisionerState

	provisionerInstance := &Provisioner{}

	// setCephCluster is only called once, on init
	// No one is allowed to modify CephCluster
	err = provisionerInstance.setCephCluster(cephClusterSpec)
	if err != nil {
		// TODO: Replace stdout to log out
		fmt.Println("[Provisioner] setCephCluster Error")
		return nil, err
	}
	err = provisionerInstance.setClientSet(clientSet)
	if err != nil {
		// TODO: Replace stdout to log out
		fmt.Println("[Provisioner] setCephClientSet Error")
		return nil, err
	}
	err = provisionerInstance.setCephNamespace(cephNamespace)
	if err != nil {
		// TODO: Replace stdout to log out
		fmt.Println("[Provisioner] setCephNamespace Error")
		return nil, err
	}
	err = provisionerInstance.setCephName(cephName)
	if err != nil {
		// TODO: Replace stdout to log out
		fmt.Println("[Provisioner] setCephName Error")
		return nil, err
	}

	pathConfigDir := provisionerInstance.getPathConfigDir()

	err = wrapper.OsWrapper.MkdirAll(pathConfigDir, 0644)
	if err != nil {
		fmt.Println("[Provisioner] CephName Driectory Create Error")
		return nil, err
	}

	// identifyProvisionerState is only called once, on init
	// No one is allowed to modify provisionerState
	currentState, err = provisionerInstance.identifyProvisionerState()
	if err != nil {
		// TODO: Replace stdout to log out
		fmt.Println("[Provisioner] identifyProvisionerState Error")
		return nil, err
	}

	err = provisionerInstance.setState(currentState)
	if err != nil {
		// TODO: Replace stdout to log out
		fmt.Println("[Provisioner] setState Error")
		return nil, err
	}

	return provisionerInstance, nil
}
