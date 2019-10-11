/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package integration

import (
	"context"
	"reflect"

	integrationv1alpha1 "github.com/wso2/k8s-ei-operator/pkg/apis/integration/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_integration")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Integration Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileIntegration{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("integration-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Integration
	err = c.Watch(&source.Kind{Type: &integrationv1alpha1.Integration{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Integration
	// Watch for deployment
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &integrationv1alpha1.Integration{},
	})

	// Watch for service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &integrationv1alpha1.Integration{},
	})

	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileIntegration{}

// ReconcileIntegration reconciles a Integration object
type ReconcileIntegration struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Integration object and makes changes based on the state read
// and what is in the Integration.Spec
// Controller logic written for creates an Integration Deployment for each Integration CR
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result. Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileIntegration) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Integration")

	// Fetch the Integration integration
	integration := &integrationv1alpha1.Integration{}
	err := r.client.Get(context.TODO(), request.NamespacedName, integration)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	deploymentObj := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: nameForDeployment(integration), Namespace: integration.Namespace}, deploymentObj)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		deployment := r.deploymentForIntegration(integration)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// Ensure the deployment replicas is the same as the spec
	replicas := integration.Spec.Replicas
	if *deploymentObj.Spec.Replicas != replicas {
		deploymentObj.Spec.Replicas = &replicas
		err = r.client.Update(context.TODO(), deploymentObj)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", deploymentObj.Namespace, "Deployment.Name", deploymentObj.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// Check if the service already exists, if not create a new one
	serviceObj := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: nameForService(integration), Namespace: integration.Namespace}, serviceObj)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		service := r.serviceForIntegration(integration)
		reqLogger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return reconcile.Result{}, err
		}
		// Service created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}

	// Update status.Status if needed
	availableReplicas := deploymentObj.Status.AvailableReplicas
	currentStatus := "NotRunning"
	if availableReplicas > 0 {
		currentStatus = "Running"
	}
	if !reflect.DeepEqual(currentStatus, integration.Status.Readiness) {
		integration.Status.Readiness = currentStatus
		err := r.client.Status().Update(context.TODO(), integration)
		if err != nil {
			reqLogger.Error(err, "Failed to update Integration status")
			return reconcile.Result{}, err
		}
	}

	// Update status.ServiceName if needed
	serviceName := nameForService(integration)
	if !reflect.DeepEqual(serviceName, integration.Status.ServiceName) {
		integration.Status.ServiceName = serviceName
		err := r.client.Status().Update(context.TODO(), integration)
		if err != nil {
			reqLogger.Error(err, "Failed to update Integration status")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// deploymentForIntegration returns a integration Deployment object
func (r *ReconcileIntegration) deploymentForIntegration(m *integrationv1alpha1.Integration) *appsv1.Deployment {
	labels := labelsForIntegration(m.Name)
	replicas := m.Spec.Replicas

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      nameForDeployment(m),
			Namespace: m.Namespace,
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
						corev1.Container{
							Image: m.Spec.Image,
							Name:  "micro-integrator",
							Ports: []corev1.ContainerPort{{
								ContainerPort: 8290,
							}},
							Env:             m.Spec.Env,
							ImagePullPolicy: corev1.PullAlways,
						},
					},
				},
			},
		},
	}
	// Set Integration instance as the owner and controller
	controllerutil.SetControllerReference(m, deployment, r.scheme)
	return deployment
}

// serviceForIntegration returns a service object
func (r *ReconcileIntegration) serviceForIntegration(m *integrationv1alpha1.Integration) *corev1.Service {
	labels := labelsForIntegration(m.Name)

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      nameForService(m),
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Port:       m.Spec.Port,
				TargetPort: intstr.FromInt(8290),
			}},
		},
	}
	// Set Integration instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// labelsForIntegration returns the labels for selecting the resources
// belonging to the given integration CR name.
func labelsForIntegration(name string) map[string]string {
	return map[string]string{"app": "integration", "integration_cr": name}
}

func nameForDeployment(m *integrationv1alpha1.Integration) string {
	return m.Name + "-deployment"
}

func nameForService(m *integrationv1alpha1.Integration) string {
	return m.Name + "-service"
}
