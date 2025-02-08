/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	scalev1 "lab-k8s-operator-sdk/api/v1"

	appsv1 "k8s.io/api/apps/v1"
)

// ScaleDeploymentReconciler reconciles a ScaleDeployment object
type ScaleDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=scale.madrid.operator.com,resources=scaledeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=scale.madrid.operator.com,resources=scaledeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=scale.madrid.operator.com,resources=scaledeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ScaleDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ScaleDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	scaleDApi := &scalev1.ScaleDeployment{}
	deploymentApi := &appsv1.Deployment{}

	err := r.Get(ctx, req.NamespacedName, scaleDApi)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("ScaleDeployment resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get scaleDeployment")
		return ctrl.Result{}, err
	} else {
		log.Info("FOUND", "scaleDeployment", scaleDApi.Name)
	}

	err = r.Get(ctx, types.NamespacedName{
		Name:      scaleDApi.Spec.DeploymentName,
		Namespace: scaleDApi.Namespace,
	}, deploymentApi)

	if err != nil && errors.IsNotFound(err) {
		fmt.Printf("ERROR   Deployment %s NOT FOUND\n", deploymentApi.Name)
		return ctrl.Result{RequeueAfter: time.Duration(20000) * time.Millisecond}, nil
	}

	fmt.Printf("Deployment %s found with %d replicas\n", deploymentApi.Name, *deploymentApi.Spec.Replicas)

	if *deploymentApi.Spec.Replicas != scaleDApi.Spec.Replicas {
		deploymentApi.Spec.Replicas = &scaleDApi.Spec.Replicas
		r.Update(ctx, deploymentApi)
		log.Info("SCALED DEPLOYMENT REPLICAS")
	} else {
		log.Info("ALREADY SCALED")
	}

	return ctrl.Result{RequeueAfter: time.Duration(5000) * time.Millisecond}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScaleDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalev1.ScaleDeployment{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
