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
	"reflect"

	// "os/exec"

	// "strings"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	itav1alpha1 "github.com/exastro-suite/it-automation-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Instance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// _ = r.Log.WithValues("instance", req.NamespacedName)
	log := r.Log.WithValues("instance", req.NamespacedName)

	fmt.Printf("\n-------------------->\n")
	// your logic here

	instance := &itav1alpha1.Instance{}
	// fetch the Instance
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// _.Info("Instance resource not found. Ignoring since object must be deleted")
			log.Info("Instance resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Instance")
		return ctrl.Result{}, err
	}

	// デプロイメントがすでに存在するかどうかを確認し、存在しない場合は新しいデプロイメントを作成
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForInstance(instance)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}
	// サービスがすでに存在するかどうかを確認し、存在しない場合は新しいサービスを作成
	foundSvc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundSvc)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		svc := r.serviceForInstance(instance)
		log.Info("Creating a new Service", "Service.Namespace", instance.Namespace, "Service.Name", instance.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		}
	}

	// デプロイメントサイズが仕様と同じであることを確認
	size := instance.Spec.Replicas
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Instanceステータスをポッド名で更新
	// このinstanceのデプロイのポッドを一覧表示
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(labelsForInstance(instance.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.Status().Update(ctx, instance)
		if err != nil {
			log.Error(err, "Failed to update Instance status")
			return ctrl.Result{}, err
		}
	}

	fmt.Printf("\n<--------------------\n")
	return ctrl.Result{}, nil
}

func labelsForInstance(name string) map[string]string {
	return map[string]string{"app": name, "Instance_cr": name}
	// return map[string]string{"app": "Instance", "Instance_cr": name}
}

// 渡されたポッドの配列のポッド名を返却
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func int64Ptr(i int64) *int64 { return &i }
func boolPtr(b bool) *bool    { return &b }

// DeploymentForInstanceはインスタンスDeploymentオブジェクトを返却
func (r *InstanceReconciler) deploymentForInstance(m *itav1alpha1.Instance) *appsv1.Deployment {
	ls := labelsForInstance(m.Name)
	replicas := m.Spec.Replicas

	// tmp := corev1
	// fmt.Printf("%+v", tmp)
	// cmdlist := []string{
	// 	"ls -l",
	// }
	// fmt.Println(strings.Join(cmdlist, ";"))
	fmt.Println("------")
	configMaplist := &corev1.ConfigMapList{}
	fmt.Println(configMaplist)
	// getConfigMap()
	fmt.Println("------")

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:    "ita001-init",
							Image:   fmt.Sprintf("exastro/it-automation:%s", m.Spec.ReleasedVersion),
							Command: []string{"/bin/sh"},
							Args: []string{
								"-c",
								"if [ -d /var/lib/mysql ]; then chown -R mysql. /var/lib/mysql; fi",
								// "systemctl restart mariadb",
								// "while [ ! -e /var/lib/mysql/mysql.sock ]; do sleep 1; done",
								// "mysql -uroot -hlocalhost --password=ita_root_password -e 'create database foobardb;'",
								// "touch /var/lib/mysql/initialized.txt",
							},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "mysql-persistent-storage",
								MountPath: "/var/lib/mysql",
							}},
						},
					},
					Containers: []corev1.Container{{
						Name:  "ita001",
						Image: fmt.Sprintf("exastro/it-automation:%s", m.Spec.ReleasedVersion),
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: m.Spec.DbRootPassword,
							},
							{
								Name:  "MYSQL_USER",
								Value: m.Spec.DbUser,
							},
							{
								Name:  "MYSQL_PASSWORD",
								Value: m.Spec.DbPassword,
							},
							// {
							// 	Name:  "MYSQL_ALLOW_EMPTY_PASSWORD",
							// 	Value: "true",
							// },
						},
						// Command: []string{"/bin/sh", "-c", "echo", "hello world."},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
								Name:          "http",
							},
							{
								ContainerPort: 443,
								Name:          "https",
							},
							{
								ContainerPort: 3306,
								Name:          "mysql",
							},
						},
						// SecurityContext: &corev1.SecurityContext{
						// 	// RunAsUser: int64Ptr(0),
						// 	Privileged: boolPtr(true),
						// },
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "mysql-persistent-storage",
							MountPath: "/var/lib/mysql",
						}},
					}},
					// SecurityContext: &corev1.PodSecurityContext{
					// 	// RunAsUser: int64Ptr(0),
					// 	Privileged: boolPtr(true),
					// },
					RestartPolicy: "Always",
					Volumes: []corev1.Volume{{
						Name: "mysql-persistent-storage",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: m.Spec.DbStorageName,
							},
						},
					}},
				},
			},
		},
	}
	// Set Instance instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// ServiceForInstanceはインスタンスServiceオブジェクトを返却
func (r *InstanceReconciler) serviceForInstance(m *itav1alpha1.Instance) *corev1.Service {
	ls := labelsForInstance(m.Name)
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Protocol: "TCP",
					Port:     80,
					// Port: m.Spec.Ports,
					Name: "http",
				},
				{
					Protocol: "TCP",
					Port:     443,
					Name:     "https",
				},
				{
					Protocol: "TCP",
					Port:     3306,
					Name:     "mysql",
				},
			},
		},
	}

	return svc
}

func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&itav1alpha1.Instance{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		// WithOptions(controller.Options{
		// 	MaxConcurrentReconciles: 2,
		// }).
		Complete(r)
}
