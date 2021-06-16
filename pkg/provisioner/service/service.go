package service

// https://github.com/ceph/ceph/blob/master/src/python-common/ceph/deployment/service_spec.py

// Service is struct for ceph service spec yaml
type Service struct {
	// ServiceType indicates service_type field
	ServiceType string `yaml:"service_type"`
	// ServiceId indicates service_id field
	ServiceID string `yaml:"service_id"`
	// Placement indicates placement field
	Placement Placement `yaml:"placement,omitempty"`
	// Unmanaged indicates unamaged field
	Unmanaged bool `yaml:"unmanaged,omitempty"`
	// previewed_only : maybe not used
}

// SetServiceType sets ServiceType to value
func (s *Service) SetServiceType(serviceType string) error {
	s.ServiceType = serviceType
	return nil
}

// SetServiceID sets ServiceID to value
func (s *Service) SetServiceID(serviceID string) error {
	s.ServiceID = serviceID
	return nil
}

// SetPlacement sets Placement to value
func (s *Service) SetPlacement(placement Placement) error {
	s.Placement = placement
	return nil
}

// SetUnmanaged sets Unmanaged to value
func (s *Service) SetUnmanaged(unmanaged bool) error {
	s.Unmanaged = unmanaged
	return nil
}

// GetServiceType gets value of ServiceType
func (s *Service) GetServiceType() string {
	return s.ServiceType
}

// GetServiceID gets value of ServiceID
func (s *Service) GetServiceID() string {
	return s.ServiceID
}

// GetPlacement gets value of Placement
func (s *Service) GetPlacement() Placement {
	return s.Placement
}

// GetUnmanaged gets value of Unmanaged
func (s *Service) GetUnmanaged() bool {
	return s.Unmanaged
}
