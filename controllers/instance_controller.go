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

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
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
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;delete

func (reconciler *InstanceReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	customResource := &itav1alpha1.Instance{}
	requeue, result, err := reconciler.fetchCustomResource(ctx, request, customResource)
	if requeue {
		return result, err
	}

	frontendDeploymentFactory := &DeploymentFactoryForFrontend{CustomResource: customResource, Reconciler: reconciler}
	requeue, result, err = reconciler.ensureK8sResource(ctx, request, frontendDeploymentFactory)
	if requeue {
		return result, err
	}

	frontendServiceFactory := &ServiceFactoryForFrontend{CustomResource: customResource, Reconciler: reconciler}
	requeue, result, err = reconciler.ensureK8sResource(ctx, request, frontendServiceFactory)
	if requeue {
		return result, err
	}

	databaseServiceFactory := &ServiceFactoryForDatabase{CustomResource: customResource, Reconciler: reconciler}
	requeue, result, err = reconciler.ensureK8sResource(ctx, request, databaseServiceFactory)
	if requeue {
		return result, err
	}

	return ctrl.Result{}, nil
}

func (reconciler *InstanceReconciler) fetchCustomResource(ctx context.Context, request ctrl.Request, customResource *itav1alpha1.Instance) (bool, ctrl.Result, error) {
	err := reconciler.Get(ctx, request.NamespacedName, customResource)
	if err != nil {
		if errors.IsNotFound(err) {
			reconciler.Log.Info("Custom resource is not found. Ignoring since object must be deleted", k8sResourceToLogParameters(customResource)...)
			return makeReturnValuesStop()
		}

		reconciler.Log.Error(err, "Failed to get custom resource", k8sResourceToLogParameters(customResource)...)

		return makeReturnValuesRequeueWithError(err)
	}

	reconciler.Log.Info("Custom resource is found.", k8sResourceToLogParameters(customResource)...)

	return makeReturnValuesContinue()
}

func (reconciler *InstanceReconciler) ensureK8sResource(ctx context.Context, request ctrl.Request, k8sResourceFactory K8sResourceFactory) (bool, ctrl.Result, error) {
	k8sResource := k8sResourceFactory.NewDefault()
	err := reconciler.Get(ctx, k8sResourceFactory.GetNamespaceName(), k8sResource)
	if err != nil && errors.IsNotFound(err) {
		k8sResource = k8sResourceFactory.New()

		reconciler.Log.Info("Creating resource", k8sResourceToLogParameters(k8sResource)...)

		err = reconciler.Create(ctx, k8sResource)
		if err != nil {
			reconciler.Log.Error(err, "Failed to create resource", k8sResourceToLogParameters(k8sResource)...)
			return makeReturnValuesRequeueWithError(err)
		}

		return makeReturnValuesRequeue()
	} else if err != nil {
		reconciler.Log.Error(err, "Failed to get resource", k8sResourceToLogParameters(k8sResource)...)
		return makeReturnValuesRequeueWithError(err)
	}

	reconciler.Log.Info("Resource is found", k8sResourceToLogParameters(k8sResource)...)

	return makeReturnValuesContinue()
}

func k8sResourceToLogParameters(k8sResource client.Object) []interface{} {
	return []interface{}{
		"group", k8sResource.GetObjectKind().GroupVersionKind().Group,
		"version", k8sResource.GetObjectKind().GroupVersionKind().Version,
		"kind", k8sResource.GetObjectKind().GroupVersionKind().Kind,
		"namespace", k8sResource.GetNamespace(),
		"name", k8sResource.GetName(),
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

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&itav1alpha1.Instance{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
