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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	itav1alpha1 "github.com/exastro-suite/it-automation-operator/api/v1alpha1"
)

type DeploymentFactoryForFrontend struct {
	Reconciler     *InstanceReconciler
	CustomResource *itav1alpha1.Instance
}

func (factory *DeploymentFactoryForFrontend) GetName() string {
	return factory.CustomResource.Name + "-frontend"
}

func (factory *DeploymentFactoryForFrontend) GetNamespace() string {
	return factory.CustomResource.Namespace
}

func (factory *DeploymentFactoryForFrontend) GetNamespaceName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: factory.GetNamespace(),
		Name:      factory.GetName(),
	}
}

func (factory *DeploymentFactoryForFrontend) NewDefault() client.Object {
	return &appsv1.Deployment{}
}

func (factory *DeploymentFactoryForFrontend) New() client.Object {
	labels := createLabels(factory.CustomResource)
	replicas := int32(1)
	privileged := true

	k8sDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: factory.GetNamespace(),
			Name:      factory.GetName(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "it-automation",
							Image: fmt.Sprintf("ghcr.io/exastro-suite/it-automation:%s-ubi8-%s", factory.CustomResource.Spec.Version, factory.CustomResource.Spec.Language),
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 80,
								},
								{
									Name:          "https",
									ContainerPort: 443,
								},
								{
									Name:          "mysql",
									ContainerPort: 3306,
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: &privileged,
							},
							/*
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "database-volume",
										MountPath: "/var/lib/mysql",
									},
								},
							*/
						},
					},
					RestartPolicy: "Always",
					/*
						Volumes: []corev1.Volume{
							{
								Name: "database-volume",
								VolumeSource: corev1.VolumeSource{
									PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
										ClaimName: factory.CustomResource.Spec.DatabasePvcName,
									},
								},
							},
						},
					*/
				},
			},
		},
	}

	// Set resource as the owner and controller
	ctrl.SetControllerReference(factory.CustomResource, k8sDeployment, factory.Reconciler.Scheme)

	return k8sDeployment
}
