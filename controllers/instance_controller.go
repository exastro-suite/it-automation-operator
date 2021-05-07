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
	"context"
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	itav1alpha1 "github.com/exastro-suite/it-automation-operator/api/v1alpha1"
)

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ita.cr.exastro,resources=instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ita.cr.exastro,resources=instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ita.cr.exastro,resources=instances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Instance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (reconciler *InstanceReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	customResource := &itav1alpha1.Instance{}

	yield, result, err := reconciler.fetchCustomResource(ctx, request, customResource)
	if yield {
		return result, err
	}

	yield, result, err = reconciler.ensureK8sDeployment(ctx, request, customResource)
	if yield {
		return result, err
	}

	yield, result, err = reconciler.ensureK8sService(ctx, request, customResource)
	if yield {
		return result, err
	}

	return ctrl.Result{}, nil
}

func createLabels(customResource *itav1alpha1.Instance) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":     "it_automation",
		"app.kubernetes.io/instance": customResource.Name,
	}
}

func makeReturnValuesStop() (bool, ctrl.Result, error) {
	return true, ctrl.Result{}, nil
}

func makeReturnValuesRequeue() (bool, ctrl.Result, error) {
	return true, ctrl.Result{Requeue: true}, nil
}

func makeReturnValuesRequeueWithError(err error) (bool, ctrl.Result, error) {
	return true, ctrl.Result{}, err
}

func makeReturnValuesContinue() (bool, ctrl.Result, error) {
	return false, ctrl.Result{}, nil
}

func (reconciler *InstanceReconciler) fetchCustomResource(ctx context.Context, request ctrl.Request, customResource *itav1alpha1.Instance) (bool, ctrl.Result, error) {
	err := reconciler.Get(ctx, request.NamespacedName, customResource)
	if err != nil {
		if errors.IsNotFound(err) {
			reconciler.Log.Info("Custom resource not found. Ignoring since object must be deleted")
			return makeReturnValuesStop()
		}

		reconciler.Log.Error(err, "Failed to get custom resource")
		return makeReturnValuesRequeueWithError(err)
	}

	return makeReturnValuesContinue()
}

func (reconciler *InstanceReconciler) ensureK8sDeployment(ctx context.Context, request ctrl.Request, customResource *itav1alpha1.Instance) (bool, ctrl.Result, error) {
	k8sDeployment := &appsv1.Deployment{}
	err := reconciler.Get(ctx, types.NamespacedName{Name: customResource.Name, Namespace: customResource.Namespace}, k8sDeployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		k8sDeployment = reconciler.createK8sDeployment(customResource)

		reconciler.Log.Info("Creating a new Deployment", "Deployment.Namespace", k8sDeployment.Namespace, "Deployment.Name", k8sDeployment.Name)

		err = reconciler.Create(ctx, k8sDeployment)
		if err != nil {
			reconciler.Log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", k8sDeployment.Namespace, "Deployment.Name", k8sDeployment.Name)
			return makeReturnValuesRequeueWithError(err)
		}

		// Deployment created successfully - return and requeue
		return makeReturnValuesRequeue()
	} else if err != nil {
		reconciler.Log.Error(err, "Failed to get Deployment")
		return makeReturnValuesRequeueWithError(err)
	}

	// Deployment already exists
	reconciler.Log.Info("Deployment already exists", "Deployment.Namespace", k8sDeployment.Namespace, "Deployment.Name", k8sDeployment.Name)

	return makeReturnValuesContinue()
}

func (reconciler *InstanceReconciler) createK8sDeployment(customResource *itav1alpha1.Instance) *appsv1.Deployment {
	labels := createLabels(customResource)
	replicas := int32(1)
	privileged := true

	k8sDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      customResource.Name,
			Namespace: customResource.Namespace,
			Labels:    labels,
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
							Image: fmt.Sprintf("exastro/it-automation:%s", customResource.Spec.Version),
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
										ClaimName: customResource.Spec.DatabasePvcName,
									},
								},
							},
						},
					*/
				},
			},
		},
	}

	// Set Instance instance as the owner and controller
	ctrl.SetControllerReference(customResource, k8sDeployment, reconciler.Scheme)

	return k8sDeployment
}

func (reconciler *InstanceReconciler) ensureK8sService(ctx context.Context, request ctrl.Request, customResource *itav1alpha1.Instance) (bool, ctrl.Result, error) {
	k8sService := &corev1.Service{}
	err := reconciler.Get(ctx, types.NamespacedName{Name: customResource.Name, Namespace: customResource.Namespace}, k8sService)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		k8sService = reconciler.createK8sService(customResource)

		reconciler.Log.Info("Creating a new Service", "Service.Namespace", k8sService.Namespace, "Service.Name", k8sService.Name)

		err = reconciler.Create(ctx, k8sService)
		if err != nil {
			reconciler.Log.Error(err, "Failed to create new Service", "Service.Namespace", k8sService.Namespace, "Service.Name", k8sService.Name)
			return makeReturnValuesRequeueWithError(err)
		}

		// Service created successfully - return and requeue
		return makeReturnValuesRequeue()
	} else if err != nil {
		reconciler.Log.Error(err, "Failed to get Deployment")
		return makeReturnValuesRequeueWithError(err)
	}

	// Service already exists
	reconciler.Log.Info("Service already exists", "Service.Namespace", k8sService.Namespace, "Service.Name", k8sService.Name)

	return makeReturnValuesContinue()
}

func (reconciler *InstanceReconciler) createK8sService(customResource *itav1alpha1.Instance) *corev1.Service {
	labels := createLabels(customResource)

	k8sService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      customResource.Name,
			Namespace: customResource.Namespace,
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
				{
					Name:       "mysql",
					Port:       3306,
					TargetPort: intstr.FromInt(3306),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	// Set Instance instance as the owner and controller
	ctrl.SetControllerReference(customResource, k8sService, reconciler.Scheme)

	return k8sService
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&itav1alpha1.Instance{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
