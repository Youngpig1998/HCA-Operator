/*
Copyright 2022.

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
	"github.com/Youngpig1998/HCA-Operator/iaw-shared-helpers/pkg/bootstrap"
	"github.com/Youngpig1998/HCA-Operator/internal/operator"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	autoscalev1beta1 "github.com/Youngpig1998/HCA-Operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	controllerManagerName = "HCA-Operator-controller-manager"
)

// HCAJobReconciler reconciles a HCAJob object
type HCAJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
	Config *rest.Config
}

//+kubebuilder:rbac:groups=autoscale.njtech.edu.cn,resources=hcajobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=autoscale.njtech.edu.cn,resources=hcajobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=autoscale.njtech.edu.cn,resources=hcajobs/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.oam.dev,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HCAJob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *HCAJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("HCAJob", req.NamespacedName)
	//log := log.FromContext(ctx)

	log.Info("1. start reconcile logic")
	// Instantialize the data structure
	instance := &autoscalev1beta1.HCAJob{}

	//First,query the webhook instance
	err := r.Get(ctx, req.NamespacedName, instance)

	if err != nil {
		// If there is no instance, an empty result is returned, so that the Reconcile method will not be called immediately
		if errors.IsNotFound(err) {
			log.Info("Instance not found, maybe removed")
			return reconcile.Result{}, nil
		}
		log.Error(err, "query action happens error")
		// Return error message
		return ctrl.Result{}, err
	}

	//Set the bootstrapClient's owner value as the webhook,so the resources we create then will be set reference to the webhook
	//when the webhook cr is deleted,the resources(such as deployment.configmap,issuer...) we create will be deleted too
	bootstrapClient, err := bootstrap.NewClient(r.Config, r.Scheme, controllerManagerName, instance)
	if err != nil {
		log.Error(err, "failed to initialise bootstrap client")
		return ctrl.Result{}, err
	}

	//We create servicemonitors to monitor the specified microservice Services
	//serviceMonitorName := deploymentNames[i] + "-hpa"
	serviceLabels := instance.Spec.MonitorDatas.ServiceLabels
	for i := 0; i < len(serviceLabels); i++ {
		var serviceMonitorName string
		for _, value := range serviceLabels[i] {
			serviceMonitorName += value
		}
		serviceMonitorName += "-svcmonitor"
		serviceMonitor := operator.ServiceMonitor(serviceMonitorName, serviceLabels[i], instance)
		err = bootstrapClient.CreateResource(serviceMonitorName, serviceMonitor)
		if err != nil {
			log.Error(err, "failed to create "+serviceMonitorName, "Name", serviceMonitorName)
			return ctrl.Result{}, err
		}
	}

	//We create hpas to auto-scale the specified microservice Deployments
	deploymentNames := instance.Spec.ScaleDatas.ScaleTargetDeploymentNames
	for i := 0; i < len(deploymentNames); i++ {
		hpaName := deploymentNames[i] + "-hpa"
		hpa := operator.HorizontalPodAutoscaler(hpaName, instance)
		err = bootstrapClient.CreateResource(hpaName, hpa)
		if err != nil {
			log.Error(err, "failed to create "+hpaName, "Name", hpaName)
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HCAJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalev1beta1.HCAJob{}).
		Complete(r)
}
