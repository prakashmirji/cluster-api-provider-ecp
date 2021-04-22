/*


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
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrastructurev1alpha3 "github.com/prakashmirji/cluster-api-provider-ecp/api/v1alpha3"
	"github.com/prakashmirji/cluster-api-provider-ecp/internal/pkg/machine"
)

const (
	ecpMachineFinalizer = "ecpMachine.infrastructure.cluster.x-k8s.io"
	kubeClusterLabel    = "cluster"
	tagClusterName      = "clusterName"
)

// ECPMachineReconciler reconciles a ECPMachine object
type ECPMachineReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ecpmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ecpmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

func (r *ECPMachineReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("ecpmachine", req.NamespacedName)

	// Firstly fetch the custom resource ecpMachin details
	var ecpMachine infrastructurev1alpha3.ECPMachine
	if err := r.Get(ctx, req.NamespacedName, &ecpMachine); err != nil {
		log.Error(err, "unable to fetch ecpMachine")
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// examine DeletionTimestamp to determine if object is under deletion
	if !ecpMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("going to execute reconcileDelete()")
		// The object is not being deleted
		return r.reconcileDelete(&ecpMachine)
	}
	// The object is being non deleted
	return r.reconcileNormal(&ecpMachine)
}

func (r *ECPMachineReconciler) reconcileNormal(ecpMachine *infrastructurev1alpha3.ECPMachine) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("ecpmachine", ecpMachine.Name)

	secret := make(map[string]string)
	config := make(map[string]string)
	secret["dbNamespace"] = "static-worker"
	provider := "static"

	macProvider, err := newMachineProvider(provider, secret, config)
	if err != nil {
		msg := "unable to create machine provider client"
		log.Error(err, msg)
		//r.setMachineError(kubeMachine, infrav1alpha2.MachineProviderError, msg)
		return ctrl.Result{}, err
	}

	if ecpMachine.Status.Status == "" {
		ecpMachine.Status.Status = "initialized"
	}

	ecpMachine.Status.State = "pending"
	ecpMachine.Status.Ready = false

	if err := r.Status().Update(ctx, ecpMachine); err != nil {
		log.Info("error updating status at the start in reconcile normal", "name", ecpMachine.Name)
		return ctrl.Result{}, err
	}

	tags := map[string]string{}
	cluster := ecpMachine.Labels[kubeClusterLabel]
	tags[tagClusterName] = cluster

	var mac machine.Machine

	if mac, err = macProvider.GetMachine(ecpMachine.Name); err != nil {
		return ctrl.Result{}, err
	} else if mac == nil {
		// Machine does not exist, proceed to create machine
		log.Info("creating machine", "machine", ecpMachine.Name)
		// fetch default machine networks from config map
		//mapOfRoleToDefaultMachineNetworks, err := getDefaultMachineNetworks(r.Client)
		// if err != nil {
		// 	log.Error(err, "unable to fetch default machine networks configmap")
		// 	return ctrl.Result{}, err
		// }
		mapOfRoleToDefaultMachineNetworks := make(map[string]string)
		mac, err = r.createKubeMachine(macProvider, mapOfRoleToDefaultMachineNetworks, ecpMachine, tags)
		if err != nil {
			log.Error(err, "unable to create machine")
			return ctrl.Result{}, err
		}
		log.Info("successfully created machine", "machine", mac.Name(), "id", mac.ID(), "hostname", mac.Hostname())
	}

	// If the object does not have our finalizer add the finalizer and update the object. This is equivalent
	// registering our finalizer.
	if !containsString(ecpMachine.ObjectMeta.Finalizers, ecpMachineFinalizer) {
		ecpMachine.ObjectMeta.Finalizers = append(ecpMachine.ObjectMeta.Finalizers, ecpMachineFinalizer)
		if err := r.Update(ctx, ecpMachine); err != nil {
			return ctrl.Result{}, err
		}
	}

	ecpMachine.Status.Ready = true
	ecpMachine.Status.State = "created"
	ecpMachine.Status.Status = "deployed"
	ecpMachine.Status.HostIP = mac.Hostname()

	if err := r.Status().Update(ctx, ecpMachine); err != nil {
		log.Info("error updating status at the end in reconcile normal", "name", ecpMachine.Name)
		return ctrl.Result{}, err
	}
	log.Info("completed periodic machine status update")

	return ctrl.Result{Requeue: true, RequeueAfter: time.Duration(20) * time.Second}, nil

}

func (r *ECPMachineReconciler) reconcileDelete(ecpMachine *infrastructurev1alpha3.ECPMachine) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("ecpmachine", ecpMachine.Name)

	secret := make(map[string]string)
	config := make(map[string]string)
	secret["dbNamespace"] = "static-worker"
	provider := "static"

	macProvider, err := newMachineProvider(provider, secret, config)
	if err != nil {
		msg := "unable to create machine provider client"
		log.Error(err, msg)
		//r.setMachineError(kubeMachine, infrav1alpha2.MachineProviderError, msg)
		return ctrl.Result{}, err
	}

	err = macProvider.DeleteMachine(ecpMachine.Name)
	if err != nil {
		log.Error(err, "unable to delete machine from the provider")
		return ctrl.Result{}, err
	}

	ecpMachine.Status.State = "deleting"

	if err := r.Status().Update(context.Background(), ecpMachine); err != nil {
		log.Info("error updating status at the end in reconcile delete", "name", ecpMachine.Name)
		return ctrl.Result{}, err
	}

	// remove our finalizer from the list and update it.
	if containsString(ecpMachine.ObjectMeta.Finalizers, ecpMachineFinalizer) {
		ecpMachine.ObjectMeta.Finalizers = removeString(ecpMachine.ObjectMeta.Finalizers, ecpMachineFinalizer)
		if err := r.Update(ctx, ecpMachine); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// createKubeMachine create's a Machine provider machine as per supplied spec
func (r *ECPMachineReconciler) createKubeMachine(macProvider machine.Provider, mapOfRoleToDefaultMachineNetworks map[string]string, ecpMachine *infrastructurev1alpha3.ECPMachine, tags map[string]string) (machine.Machine, error) {
	log := r.Log.WithValues("ecpmachine", ecpMachine.Name)
	// r.setMachineState(kubeMachine, k8s.StateCreating)
	// r.setMachineHealth(kubeMachine, k8s.HealthUnknown)
	var networks []string
	createInfo := machine.MachineCreateData{
		OsImage:   ecpMachine.Spec.OsImage,
		OsVersion: ecpMachine.Spec.OsVersion,
		Size:      ecpMachine.Spec.Size,
		SSHUser:   ecpMachine.Spec.SSHUser,
		SSHKey:    ecpMachine.Spec.SSHKey,
		Networks:  networks,
		Proxy:     ecpMachine.Spec.Proxy,
		Tags:      tags,
	}
	log.Info("machine create started")
	return macProvider.CreateMachine(ecpMachine.Name, createInfo)
}

func (r *ECPMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha3.ECPMachine{}).
		Complete(r)
}
