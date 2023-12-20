/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// LobbySpec defines the desired state of Lobby
type LobbySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Lobby. Edit lobby_types.go to remove/update
	Image         string `json:"image"`
	Replicas      int    `json:"replicas"`
	EnableIngress bool   `json:"enable_ingress,omitempty"`
	EnableService bool   `json:"enable_service"`
	Domain        string `json:"domain"`
	Token         string `json:"token"`
	Dedicated     string `json:"dedicated,omitempty"`
}

// LobbyStatus defines the observed state of Lobby
type LobbyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Image     string `json:"image"`
	Replicas  int    `json:"replicas"`
	Domain    string `json:"domain"`
	Dedicated string `json:"dedicated"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Lobby is the Schema for the lobbies API
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.metadata.name`
// +kubebuilder:printcolumn:name="Image",type=string,JSONPath=`.spec.Image`
// +kubebuilder:printcolumn:name="Domain",type=string,JSONPath=`.spec.Domain`
// +kubebuilder:printcolumn:name="Dedicated",type=string,JSONPath=`.spec.Dedicated`
// +kubebuilder:subresource:status
type Lobby struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LobbySpec   `json:"spec,omitempty"`
	Status LobbyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LobbyList contains a list of Lobby
type LobbyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Lobby `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Lobby{}, &LobbyList{})
}
