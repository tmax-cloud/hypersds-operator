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

type CephClusterMonSpec struct {
	Count int `json:"count"`
}

type Node struct {
	Name       string         `json:"name"`
	AccessInfo NodeAccessInfo `json:"accessInfo"`
	Devices    []string       `json:"devices"`
}

type NodeAccessInfo struct {
	IP       string `json:"ip"`
	UserID   string `json:"userId"`
	Password string `json:"password"`
	HostName string `json:"hostName"`
}

// CephClusterSpec defines the desired state of CephCluster
type CephClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Mon    CephClusterMonSpec `json:"mon"`
	Nodes  []Node             `json:"nodes"`
	Config map[string]string  `json:"config,omitempty"`
}

type CephClusterState string

// CephClusterStatus defines the observed state of CephCluster
type CephClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	State CephClusterState `json:"state"`
}

// +kubebuilder:object:root=true

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
