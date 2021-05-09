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

type ServiceFactoryForFrontend struct {
	Reconciler     *InstanceReconciler
	CustomResource *itav1alpha1.Instance
}

func (factory *ServiceFactoryForFrontend) GetName() string {
	return factory.CustomResource.Name + "-frontend"
}

func (factory *ServiceFactoryForFrontend) GetNamespace() string {
	return factory.CustomResource.Namespace
}

func (factory *ServiceFactoryForFrontend) GetNamespaceName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: factory.GetNamespace(),
		Name:      factory.GetName(),
	}
}

func (factory *ServiceFactoryForFrontend) NewDefault() client.Object {
	return &corev1.Service{}
}

func (factory *ServiceFactoryForFrontend) New() client.Object {
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
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
				{
					Name:       "https",
					Port:       443,
					TargetPort: intstr.FromInt(443),
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	// Set resource as the owner and controller
	ctrl.SetControllerReference(factory.CustomResource, k8sService, factory.Reconciler.Scheme)

	return k8sService
}
