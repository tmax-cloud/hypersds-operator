package osd

import (
	"bytes"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/service"
)

type OsdSpec struct {
	DataDevices     Device   `yaml:"data_devices,omitempty"`
	DbDevices       Device   `yaml:"db_devices,omitempty"`
	WalDevices      Device   `yaml:"wal_devices,omitempty"`
	JournalDevices  Device   `yaml:"journal_devices,omitempty"`
	DataDirectories []string `yaml:"data_directories,omitempty"`
	OsdsPerDevice   int      `yaml:"osds_per_device,omitempty"`
	Objectstore     string   `yaml:"objectstore,omitempty"`
	Encrypted       bool     `yaml:"encrypted,omitempty"`
	Filter_logic    string   `yaml:"filter_logic,omitempty"`
	/*
	    				db_slots=None,  # type: Optional[int]
	                    wal_slots=None,  # type: Optional[int]
	                    osd_id_claims=None,  # type: Optional[Dict[str, List[str]]]
	                    block_db_size=None,  # type: Union[int, str, None]
	                    block_wal_size=None,  # type: Union[int, str, None]
	                    journal_size=None,  # type: Union[int, str, None]
	                    service_type=None,  # type: Optional[str]
	                    unmanaged=False,  # type: bool
	   				preview_only=False,  # type: bool
	*/
}

type Osd struct {
	Service service.Service `yaml:",inline"`
	Spec    OsdSpec         `yaml:"spec,omitempty"`
}

func (o *Osd) SetService(s service.Service) error {
	o.Service = s
	return nil
}

func (o *Osd) SetDataDevices(dataDevices Device) error {
	o.Spec.DataDevices = dataDevices
	return nil
}

func (o *Osd) GetService() service.Service {
	return o.Service
}

func (o *Osd) GetDataDevices() Device {
	return o.Spec.DataDevices
}

func (o *Osd) CompareDataDevices(targetOsd *Osd) ([]string, []string, error) {
	// o: orch osd, targetOsd: cephCr osd
	dataDevices := o.GetDataDevices()

	devicePaths := dataDevices.GetPaths()

	targetDataDevices := targetOsd.GetDataDevices()

	targetDevicePaths := targetDataDevices.GetPaths()

	var deviceMap map[string]bool
	var addDeviceList, removeDeviceList []string

	deviceMap = make(map[string]bool)

	for _, device := range devicePaths {
		deviceMap[device] = false
	}
	for _, device := range targetDevicePaths {
		_, exists := deviceMap[device]
		if exists {
			deviceMap[device] = true
		} else {
			addDeviceList = append(addDeviceList, device)
		}
	}

	for device, value := range deviceMap {
		if !value {
			removeDeviceList = append(removeDeviceList, device)
		}
	}
	return addDeviceList, removeDeviceList, nil
}

func (o *Osd) MakeYmlFile(yaml wrapper.YamlInterface, ioUtilWrapper wrapper.IoUtilInterface, fileName string) error {
	osdYaml, err := yaml.Marshal(o)
	if err != nil {
		return err
	}

	err = ioUtilWrapper.WriteFile(fileName, osdYaml, 0644)
	return err
}

func NewOsdsFromCephCr(cephSpec hypersdsv1alpha1.CephClusterSpec) ([]*Osd, error) {
	var osds []*Osd

	for _, osdSpec := range cephSpec.Osd {
		var hosts []string
		var osd Osd
		var dataDevices Device
		var s service.Service
		var placement service.Placement

		// set Placement, Service
		hosts = append(hosts, osdSpec.HostName)
		err := placement.SetHosts(hosts)
		if err != nil {
			return nil, err
		}
		err = s.SetPlacement(placement)
		if err != nil {
			return nil, err
		}
		err = s.SetServiceType("osd")
		if err != nil {
			return nil, err
		}
		err = s.SetServiceId("osd_" + osdSpec.HostName)
		if err != nil {
			return nil, err
		}
		// set device
		err = dataDevices.SetPaths(osdSpec.Devices)
		if err != nil {
			return nil, err
		}
		err = osd.SetDataDevices(dataDevices)
		if err != nil {
			return nil, err
		}
		err = osd.SetService(s)
		if err != nil {
			return nil, err
		}
		osds = append(osds, &osd)
	}

	return osds, nil
}

func NewOsdsFromCephOrch(yaml wrapper.YamlInterface, rawOsdsFromOrch []byte) ([]*Osd, error) {

	var osds []*Osd
	readerOrch := bytes.NewReader(rawOsdsFromOrch)
	dec := yaml.NewDecoder(readerOrch)
	for {
		var osd Osd
		err := dec.Decode(&osd)
		if err != nil {
			break
		}
		osds = append(osds, &osd)
	}
	return osds, nil
}
