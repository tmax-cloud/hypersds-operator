package provisioner

import (
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/util"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/config"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/node"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"context"
	"fmt"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Provisioner contains information for ceph deployment and current ceph deployment status
type Provisioner struct {
	cephCluster hypersdsv1alpha1.CephClusterSpec

	cephNamespace string
	cephName      string
	clientSet     client.Client
	cephConfig    *config.CephConfig
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
	// Create config object from Ceph CR
	err := p.setCephConfig()
	if err != nil {
		return err
	}

	nodeList, err := p.getNodes()
	if err != nil {
		return err
	}

	err = installPackages(nodeList)
	if err != nil {
		return err
	}

	deployNode, err := p.getDeployNode()
	if err != nil {
		return err
	}
	err = p.bootstrap(deployNode)
	if err != nil {
		return err
	}

	err = p.updateKubeObject(deployNode)
	if err != nil {
		return err
	}

	cephConf, cephKeyring, err := p.getCephData()
	if err != nil {
		if kubeerrors.IsNotFound(err) {
			if err2 := p.updateKubeObject(deployNode); err2 != nil {
				return err
			}
		}
		return err
	}

	err = p.applyHost(wrapper.YamlWrapper, wrapper.ExecWrapper, wrapper.IoUtilWrapper, cephConf, cephKeyring)
	if err != nil {
		return err
	}
	err = p.applyOsd(cephConf, cephKeyring)
	if err != nil {
		return err
	}

	pathConfigDir := p.getPathConfigDir()
	err = wrapper.OsWrapper.RemoveAll(pathConfigDir)
	if err != nil {
		fmt.Println("[Provisioner] CephName Directory Remove Error")
		return err
	}

	return nil
}

func (p *Provisioner) getDeployNode() (*node.Node, error) {
	deployNode, err := p.fetchDeployNodeFromK8sStatus()
	if err != nil {
		return nil, err
	}

	return deployNode, nil
}

func (p *Provisioner) fetchDeployNodeFromK8sStatus() (*node.Node, error) {
	clientSet := p.getClientSet()
	cephNamespace := p.getCephNamespace()
	cephName := p.getCephName()

	cephCluster := &hypersdsv1alpha1.CephCluster{}
	if err := clientSet.Get(context.TODO(), types.NamespacedName{Namespace: cephNamespace, Name: cephName}, cephCluster); err != nil {
		return nil, err
	}

	deployNode, err := convertNodeType(cephCluster.Status.DeployNode)
	if err != nil {
		return nil, err
	}

	return deployNode, nil
}

func convertNodeType(k8sNode hypersdsv1alpha1.Node) (*node.Node, error) {
	hostSpec := node.HostSpec{
		Addr:     k8sNode.IP,
		HostName: k8sNode.HostName,
	}

	var n node.Node
	var err error
	err = n.SetUserID(k8sNode.UserID)
	if err != nil {
		return nil, err
	}
	err = n.SetUserPw(k8sNode.Password)
	if err != nil {
		return nil, err
	}
	err = n.SetHostSpec(&hostSpec)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (p *Provisioner) bootstrap(deployNode *node.Node) error {
	bootstrapped, err := isBootstrapped(deployNode)
	if err != nil {
		return err
	}
	if bootstrapped {
		return nil
	}

	err = installCephadm(deployNode)
	if err != nil {
		return err
	}

	cephConfig := p.getCephConfig()
	pathConfigDir := p.getPathConfigDir()
	pathConfFromCr := pathConfigDir + cephConfNameFromCr

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
	return nil
}

func isBootstrapped(deployNode *node.Node) (bool, error) {
	const checkBootstrappedCmd = "cephadm shell -- ceph -s"
	output, err := deployNode.RunSSHCmd(wrapper.SSHWrapper, checkBootstrappedCmd)
	if err != nil {
		const objectNotFound = "ObjectNotFound"
		if strings.Contains(output.String(), objectNotFound) || strings.Contains(output.String(), cmdNotFound) {
			return false, nil
		}
		// Confirm some Ceph health is returned
		const cephHealthPrefix = "HEALTH_ERR"
		if strings.Contains(output.String(), cephHealthPrefix) {
			fmt.Println("[Provisioner] ceph cluster status is error")
			return false, err
		}
		fmt.Println("[Provisioner] ceph bootstrap check is failed")
		return false, err
	}
	return true, nil
}

func (p *Provisioner) updateKubeObject(deployNode *node.Node) error {
	pathConfigDir := p.getPathConfigDir()
	pathConf := pathConfigDir + util.CephConfName
	pathKeyring := pathConfigDir + util.CephKeyringName

	// Copy conf and keyring from deploy node
	err := copyFile(deployNode, node.SOURCE, pathConfFromAdm, pathConf)
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
	return nil
}

func (p *Provisioner) getCephData() (confDataBuf, keyringDataBuf []byte, err error) {
	conf, keyring, err := p.checkKubeObjectUpdated()
	if err != nil {
		fmt.Println("[Provisioner] k8s configmap and secret check is failed")
		return conf, keyring, err
	}
	if conf == nil || keyring == nil {
		fmt.Println("[Provisioner] k8s configmap and secret are not updated")
		return conf, keyring, kubeerrors.NewNotFound(v1.Resource("CephCluster"), p.cephName)
	}
	return conf, keyring, nil
}

// TODO: Replace config const to inputs (e.g. K8sConfigMap, etc)
func (p *Provisioner) checkKubeObjectUpdated() (confDataBuf, keyringDataBuf []byte, err error) {
	// Check ceph.conf is updated to ConfigMap
	configMap := &corev1.ConfigMap{}
	if err := p.clientSet.Get(context.TODO(), types.NamespacedName{Namespace: p.cephNamespace, Name: p.cephName + util.K8sConfigMapSuffix}, configMap); err != nil {
		// Configmap must exist
		if kubeerrors.IsNotFound(err) {
			// TODO: Replace stdout to log out
			fmt.Println("ConfigMap must exist")
		}
		return nil, nil, err
	}

	confData, ok := configMap.Data["conf"]
	if !ok {
		return nil, nil, nil
	}
	confDataBuf = []byte(confData)

	// Check client.admin.keyring is updated to Secret
	secret := &corev1.Secret{}
	if err := p.clientSet.Get(context.TODO(), types.NamespacedName{Namespace: p.cephNamespace, Name: p.cephName + util.K8sSecretSuffix}, secret); err != nil {
		if kubeerrors.IsNotFound(err) {
			// TODO: Replace stdout to log out
			fmt.Println("Secret must exist")
		}
		return nil, nil, err
	}

	keyringDataBuf, ok = secret.Data["keyring"]
	if !ok {
		return nil, nil, nil
	}

	return confDataBuf, keyringDataBuf, nil
}

// NewProvisioner creates Provisioner using ceph deployment information and checks the current ceph deployment status in node
func NewProvisioner(cephClusterSpec hypersdsv1alpha1.CephClusterSpec, clientSet client.Client, cephNamespace, cephName string) (*Provisioner, error) {
	var err error

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
		fmt.Println("[Provisioner] CephName Directory Create Error")
		return nil, err
	}

	return provisionerInstance, nil
}
