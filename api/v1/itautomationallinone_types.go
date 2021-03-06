/*
Copyright 2021 NEC Corporation.

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

// ITAutomationAllInOneSpec defines the desired state of ITAutomationAllInOne
type ITAutomationAllInOneSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Pattern=`^[1-9][0-9]*\.[0-9]+\.[0-9]+$`
	// +kubebuilder:validation:Required
	Version string `json:"version,omitempty"`

	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=2
	// +kubebuilder:default=en
	Language string `json:"language,omitempty"`

	// +kubebuilder:validation:Required
	FilePvcName string `json:"filePvcName,omitempty"`

	// +kubebuilder:validation:Required
	DatabasePvcName string `json:"databasePvcName,omitempty"`
}

// ITAutomationAllInOneStatus defines the observed state of ITAutomationAllInOne
type ITAutomationAllInOneStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`

// ITAutomationAllInOne is the Schema for the itautomationallinones API
type ITAutomationAllInOne struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ITAutomationAllInOneSpec   `json:"spec,omitempty"`
	Status ITAutomationAllInOneStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ITAutomationAllInOneList contains a list of ITAutomationAllInOne
type ITAutomationAllInOneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ITAutomationAllInOne `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ITAutomationAllInOne{}, &ITAutomationAllInOneList{})
}
