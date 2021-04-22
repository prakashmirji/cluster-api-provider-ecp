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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ECPMachineSpec defines the desired state of ECPMachine
type ECPMachineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ECPMachine. Edit ECPMachine_types.go to remove/update
	//Foo string `json:"foo,omitempty"`

	//MachineType indicates machine types like onprem or aws
	// +optional
	MachineType string `json:"machineType,eomitempty"`

	//Location is the DC location name where this machine present
	// +optional
	Location string `json:"location,omitempty"`

	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	// OS Image for the machine
	// +optional
	OsImage string `json:"osImage"`
	// OS version for the machine
	// +optional
	OsVersion string `json:"osVersion"`
	// Size of the machine
	// +optional
	Size string `json:"size"`
	// SSHKey to be used
	// +optional
	SSHKey string `json:"sshKey,omitempty"`
	// SSHUser to be used
	// +optional
	SSHUser string `json:"sshUser,omitempty"`
	// SSHPassword to be used
	// +optional
	SSHPassword string `json:"sshPassword,omitempty"`
	// Proxy address for external communication
	// +optional
	Proxy string `json:"proxy,omitempty"`
	// Roles to be applied to the machine. Can take - controlplane, etcd, worker
	Roles []string `json:"roles"`
	// Tags to be applied to the machine. Can take key value pair strings
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

// ECPMachineStatus defines the observed state of ECPMachine
type ECPMachineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//Names []string `json:"names"`

	// Status stores machine status
	Status string `json:"status,omitempty"`

	// HostIP ip of the machine
	HostIP string `json:"hostIP,omitemoty"`

	// State stores machine state
	State string `json:"state,omitempty"`

	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="hostip",type=string,JSONPath=`.status.hostIP`
// +kubebuilder:printcolumn:name="osimage",type=string,JSONPath=`.spec.osImage`
// +kubebuilder:printcolumn:name="providerID",type=string,JSONPath=`.spec.providerID`
// +kubebuilder:printcolumn:name="roles",type=string,JSONPath=`.spec.roles`
// +kubebuilder:printcolumn:name="tags",type=string,JSONPath=`.spec.tags`
// +kubebuilder:printcolumn:name="status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="state",type=string,JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="ready",type=boolean,JSONPath=`.status.Ready`
// +kubebuilder:subresource:status
// +kubebuilder:object:root=true

// ECPMachine is the Schema for the ecpmachines API
type ECPMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ECPMachineSpec   `json:"spec,omitempty"`
	Status ECPMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ECPMachineList contains a list of ECPMachine
type ECPMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ECPMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ECPMachine{}, &ECPMachineList{})
}
