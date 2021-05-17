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

package controllers

import (
	itaallinonev1 "github.com/exastro-suite/it-automation-operator/api/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sResourceFactory interface {
	GetName() string
	GetNamespace() string
	GetNamespaceName() types.NamespacedName
	NewDefault() client.Object
	New() client.Object
}

func createLabels(customResource *itaallinonev1.ITAutomationAllInOne) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":     "it-automation-all-in-one",
		"app.kubernetes.io/instance": customResource.Name,
	}
}
