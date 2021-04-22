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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ECPClusterSpec defines the desired state of ECPCluster
type ECPClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ECPCluster. Edit ECPCluster_types.go to remove/update
	//Foo string `json:"foo,omitempty"`
	Clustertype string `json:"clustertype,omitempty"`
	Location    string `json:"location,omitempty"`

	// ControlPlaneEndpoint identifies the endpoint used to connect to the target’s cluster
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
}

// // APIEndpoint represents a reachable Kubernetes API endpoint.
// type APIEndpoint struct {
// 	// The hostname on which the API server is serving.
// 	Host string `json:"host"`

// 	// The port on which the API server is serving.
// 	Port int32 `json:"port"`
// }

// ECPClusterStatus defines the observed state of ECPCluster
type ECPClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status string `json:"status,omitempty"`

	// a boolean field that is true when the infrastructure is ready to be used.
	// +optional
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="clustertype",type=string,JSONPath=`.spec.clustertype`
// +kubebuilder:printcolumn:name="location",type=string,JSONPath=`.spec.location`
// +kubebuilder:printcolumn:name="status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="ready",type=boolean,JSONPath=`.status.ready`
// +kubebuilder:subresource:status
// +kubebuilder:object:root=true

// ECPCluster is the Schema for the ecpclusters API
type ECPCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ECPClusterSpec   `json:"spec,omitempty"`
	Status ECPClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ECPClusterList contains a list of ECPCluster
type ECPClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ECPCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ECPCluster{}, &ECPClusterList{})
}
