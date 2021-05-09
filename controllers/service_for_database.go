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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	itav1alpha1 "github.com/exastro-suite/it-automation-operator/api/v1alpha1"
)

type ServiceFactoryForDatabase struct {
	Reconciler     *InstanceReconciler
	CustomResource *itav1alpha1.Instance
}

func (factory *ServiceFactoryForDatabase) GetName() string {
	return factory.CustomResource.Name + "-database"
}

func (factory *ServiceFactoryForDatabase) GetNamespace() string {
	return factory.CustomResource.Namespace
}

func (factory *ServiceFactoryForDatabase) GetNamespaceName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: factory.GetNamespace(),
		Name:      factory.GetName(),
	}
}

func (factory *ServiceFactoryForDatabase) NewDefault() client.Object {
	return &corev1.Service{}
}

func (factory *ServiceFactoryForDatabase) New() client.Object {
	labels := createLabels(factory.CustomResource)

	k8sService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: factory.GetNamespace(),
			Name:      factory.GetName(),
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "mysql",
					Port:       3306,
					TargetPort: intstr.FromInt(3306),
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	// Set resource as the owner and controller
	ctrl.SetControllerReference(factory.CustomResource, k8sService, factory.Reconciler.Scheme)

	return k8sService
}
