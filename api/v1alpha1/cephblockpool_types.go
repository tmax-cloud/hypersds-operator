/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:object:root=true

// CephBlockPoolList contains a list of CephBlockPool
type CephBlockPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CephBlockPool `json:"items"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CephBlockPool is the Schema for the cephblockpools API
type CephBlockPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CephBlockPoolSpec   `json:"spec,omitempty"`
	Status CephBlockPoolStatus `json:"status,omitempty"`
}

// CephBlockPoolSpec defines the desired state of CephBlockPool
type CephBlockPoolSpec struct {
	// The failure domain: osd/host/(region or zone if available) - technically also any type in the crush map
	// +optional
	FailureDomain string `json:"failureDomain,omitempty"`

	// The root of the crush hierarchy utilized by the pool
	// +optional
	// +nullable
	CrushRoot string `json:"crushRoot,omitempty"`

	// The device class the OSD should set to for use in the pool
	// +optional
	// +nullable
	DeviceClass string `json:"deviceClass,omitempty"`

	// The inline compression mode in Bluestore OSD to set to (options are: none, passive, aggressive, force)
	// +kubebuilder:validation:Enum=none;passive;aggressive;force;""
	// +kubebuilder:default=none
	// +optional
	// +nullable
	CompressionMode string `json:"compressionMode,omitempty"`

	// The replication settings
	// +optional
	Replicated ReplicatedSpec `json:"replicated,omitempty"`

	// The erasure code settings
	// +optional
	ErasureCoded ErasureCodedSpec `json:"erasureCoded,omitempty"`

	// EnableRBDStats is used to enable gathering of statistics for all RBD images in the pool
	EnableRBDStats bool `json:"enableRBDStats,omitempty"`

	// The quota settings
	// +optional
	// +nullable
	Quotas QuotaSpec `json:"quotas,omitempty"`
}

// ReplicatedSpec represents the spec for replication in a pool
type ReplicatedSpec struct {
	// Size - Number of copies per object in a replicated storage pool, including the object itself (required for replicated pool type)
	// +kubebuilder:validation:Minimum=0
	Size uint `json:"size"`

	// TargetSizeRatio gives a hint (%) to Ceph in terms of expected consumption of the total cluster capacity
	// +optional
	TargetSizeRatio resource.Quantity `json:"targetSizeRatio,omitempty"`
}

// ErasureCodedSpec represents the spec for erasure code in a pool
type ErasureCodedSpec struct {
	// Number of coding chunks per object in an erasure coded storage pool (required for erasure-coded pool type)
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=9
	CodingChunks uint `json:"codingChunks"`

	// Number of data chunks per object in an erasure coded storage pool (required for erasure-coded pool type)
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=9
	DataChunks uint `json:"dataChunks"`

	// The algorithm for erasure coding
	// +optional
	Algorithm string `json:"algorithm,omitempty"`
}

// QuotaSpec represents the spec for quotas in a pool
type QuotaSpec struct {
	// MaxSize represents the quota in bytes as a string
	// +kubebuilder:validation:Pattern=`^[0-9]+[\.]?[0-9]*([KMGTPE]i|[kMGTPE])?$`
	// +optional
	MaxSize *string `json:"maxSize,omitempty"`

	// MaxObjects represents the quota in objects
	// +optional
	MaxObjects *uint64 `json:"maxObjects,omitempty"`
}

// CephBlockPoolStatus defines the observed state of CephBlockPool
type CephBlockPoolStatus struct {
	State      CephBlockPoolState `json:"state"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// CephBlockPoolState is the current state of CephBlockPool
type CephBlockPoolState string

const (
	// CephBlockPoolConditionReadyToUse indicates CephBlockPool is ready to use
	CephBlockPoolConditionReadyToUse = "ReadyToUse"
)

func init() {
	SchemeBuilder.Register(&CephBlockPool{}, &CephBlockPoolList{})
}
