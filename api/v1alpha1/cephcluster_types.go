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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CephClusterMonSpec defines the spec for monitor related option
type CephClusterMonSpec struct {
	Count int `json:"count"`
}

// CephClusterOsdSpec defines the spec for osd related option
type CephClusterOsdSpec struct {
	HostName string   `json:"hostName"`
	Devices  []string `json:"devices"`
}

// Node defines the spec for node related option
type Node struct {
	IP       string `json:"ip"`
	UserID   string `json:"userId"`
	Password string `json:"password"`
	HostName string `json:"hostName"`
}

// CephClusterSpec defines the desired state of CephCluster
type CephClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Mon    CephClusterMonSpec   `json:"mon"`
	Osd    []CephClusterOsdSpec `json:"osd"`
	Nodes  []Node               `json:"nodes"`
	Config map[string]string    `json:"config,omitempty"`
}

// CephClusterCondition is the current condition of CephCluster
type CephClusterCondition string

const (
	// ConditionReadyToUse indicates CephCluster is ready to use
	ConditionReadyToUse CephClusterCondition = "ReadyToUse"
	// ConditionBootstrapped indicates CephCluster is bootstrapped
	ConditionBootstrapped CephClusterCondition = "Bootstrapped"
	// ConditionOsdDeployed indicates CephCluster is deployed with osds
	ConditionOsdDeployed CephClusterCondition = "OsdDeployed"
)

// CephClusterState is the current state of CephCluster
type CephClusterState string

const (
	// CephClusterStatePending indicates CephClusterState is creating or updating
	CephClusterStatePending CephClusterState = "Pending"
	// CephClusterStateCreating indicates CephClusterState is creating
	CephClusterStateCreating CephClusterState = "Creating"
	// CephClusterStateUpdating indicates CephClusterState is updating
	CephClusterStateUpdating CephClusterState = "Updating"
	// CephClusterStateRunning indicates CephClusterState is available
	CephClusterStateRunning CephClusterState = "Running"
	// CephClusterStateError indicates CephClusterState is error
	CephClusterStateError CephClusterState = "Error"
)

// CephClusterStatus defines the observed state of CephCluster
type CephClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DeployNode Node               `json:"deployNode,omitempty"`
	State      CephClusterState   `json:"state"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="Current state of CephCluster"

// CephCluster is the Schema for the cephclusters API
type CephCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CephClusterSpec   `json:"spec,omitempty"`
	Status CephClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CephClusterList contains a list of CephCluster
type CephClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CephCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CephCluster{}, &CephClusterList{})
}
