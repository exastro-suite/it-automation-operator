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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// InstanceSpec defines the desired state of Instance
type InstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=2
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Pattern=`^[1-9][0-9]?\.[0-9]\.[0-9](|-ja)$`
	ReleasedVersion string `json:"releasedVersion,omitempty"`

	// +kubebuilder:default="10.3"
	// +kubebuilder:validation:Pattern=`^[1-9][0-9]?\.[0-9]$`
	DbVersion string `json:"dbVersion,omitempty"`

	// +kubebuilder:validation:MinLength=4
	DbName string `json:"dbName,omitempty"`

	// +kubebuilder:validation:MinLength=8
	// +kubebuilder:validation:MaxLength=30
	DbRootPassword string `json:"dbRootPassword,omitempty"`

	// +kubebuilder:validation:MinLength=4
	// +kubebuilder:validation:MaxLength=30
	DbUser string `json:"dbUser,omitempty"`

	// +kubebuilder:validation:MinLength=4
	// +kubebuilder:validation:MaxLength=20
	DbPassword string `json:"dbPassword,omitempty"`

	// +kubebuilder:validation:MinLength=4
	// +kubebuilder:validation:MLength=30
	DbStorageName string `json:"dbStorageName,omitempty"`

	// +kubebuilder:validation:MinLength=4
	// +kubebuilder:validation:MaxLength=30
	ItaStorageName string `json:"itaStorageName,omitempty"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Ports int32 `json:"ports,omitempty"`
}

// InstanceStatus defines the observed state of Instance
type InstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Instance is the Schema for the instances API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Replicas",type=string,JSONPath=`.spec.replicas`
// +kubebuilder:printcolumn:name="Released",type=string,JSONPath=`.spec.releasedVersion`
// +kubebuilder:printcolumn:name="DbVer",type=string,JSONPath=`.spec.dbVersion`
// +kubebuilder:printcolumn:name="Storage0",type=string,JSONPath=`.spec.storage0`
// +kubebuilder:printcolumn:name="Storage1",type=string,JSONPath=`.spec.storage1`
// // +kubebuilder:printcolumn:name="Bravely Run Away",type=boolean,JSONPath=`.spec.targetPorts[?(@ == "Sir Robin")]`,description="when danger rears its ugly head, he bravely turned his tail and fled",priority=10
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec,omitempty"`
	Status InstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Instance{}, &InstanceList{})
}
