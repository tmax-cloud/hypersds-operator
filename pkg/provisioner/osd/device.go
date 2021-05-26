package osd

// https://github.com/ceph/ceph/blob/master/src/python-common/ceph/deployment/drive_group.py

// Device is struct to define ceph device spec of ceph spec yaml
type Device struct {
	// Paths indicates paths field
	Paths []string `yaml:"paths,omitempty"`
	// Model indicates model field
	Model string `yaml:"model,omitempty"`
	// Size indicates size field
	Size string `yaml:"size,omitempty"`
	// Rotational indicates rotational field
	Rotational bool `yaml:"rotational,omitempty"`
	// Limit indicates limit field
	Limit int `yaml:"limit,omitempty"`
	// Vendor indicates vendor field
	Vendor string `yaml:"vendor,omitempty"`
	// All indicates all field
	All bool `yaml:"all,omitempty"`
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (d *Device) setPaths(paths []string) error {
	d.Paths = paths
	return nil
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (d *Device) setModel(model string) error {
	d.Model = model
	return nil
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (d *Device) setSize(size string) error {
	d.Size = size
	return nil
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (d *Device) setRotational(rotational bool) error {
	d.Rotational = rotational
	return nil
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (d *Device) setLimit(limit int) error {
	d.Limit = limit
	return nil
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (d *Device) setVendor(vendor string) error {
	d.Vendor = vendor
	return nil
}

//nolint:unparam // now, setMethod always return nil, if setMethod needs return other value, deletes nolint and implement
func (d *Device) setAll(all bool) error {
	d.All = all
	return nil
}

func (d *Device) getPaths() []string {
	return d.Paths
}
func (d *Device) getModel() string {
	return d.Model
}
func (d *Device) getSize() string {
	return d.Size
}
func (d *Device) getRotational() bool {
	return d.Rotational
}
func (d *Device) getLimit() int {
	return d.Limit
}
func (d *Device) getVendor() string {
	return d.Vendor
}
func (d *Device) getAll() bool {
	return d.All
}
