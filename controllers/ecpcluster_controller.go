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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrastructurev1alpha3 "github.com/prakashmirji/cluster-api-provider-ecp/api/v1alpha3"
)

// ECPClusterReconciler reconciles a ECPCluster object
type ECPClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ecpclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ecpclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

func (r *ECPClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("ecpcluster", req.NamespacedName)

	var cluster infrastructurev1alpha3.ECPCluster

	if err := r.Get(ctx, req.NamespacedName, &cluster); err != nil {
		log.Info("error getting object", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if cluster.Status.Status == "OK" {
		return ctrl.Result{}, nil
	}

	if cluster.Status.Status == "NOT_SUPPORTED" {
		return ctrl.Result{}, nil
	}

	ctype := cluster.Spec.Clustertype
	switch ctype {
	case "ecp":
		cluster.Status.Status = "OK"
	case "aws":
		cluster.Status.Status = "NOT_SUPPORTED"
	}

	cluster.Status.Ready = true
	cluster.Spec.ControlPlaneEndpoint = clusterv1.APIEndpoint{
		Host: "10.10.10.10",
		Port: 8080,
	}

	if err := r.Status().Update(ctx, &cluster); err != nil {
		log.Info("error updating status", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("hello,  this is from clustertype controller", "name", req.NamespacedName)

	return ctrl.Result{}, nil
}

func (r *ECPClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha3.ECPCluster{}).
		Complete(r)
}
