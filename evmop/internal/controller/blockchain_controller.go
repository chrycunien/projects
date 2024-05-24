/*
Copyright 2024.

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
	"log"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	learnv1alpha1 "evmop/api/v1alpha1"
)

// BlockchainReconciler reconciles a Blockchain object
type BlockchainReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=learn.gocrazy.com,resources=blockchains,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=learn.gocrazy.com,resources=blockchains/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=learn.gocrazy.com,resources=blockchains/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Blockchain object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *BlockchainReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.SetPrefix("BlockchainReconciler")

	blockchain := &learnv1alpha1.Blockchain{}
	err := r.Get(ctx, req.NamespacedName, blockchain)
	if err != nil {
		return ctrl.Result{}, nil
	}

	log.Println("namespace", blockchain.Namespace, blockchain.GetNamespace(), req.NamespacedName)
	log.Println("name", blockchain.Name)
	log.Println("replicas", *blockchain.Spec.Replicas)
	log.Println("image", blockchain.Spec.Image)

	for _, value := range blockchain.Spec.Command {
		log.Printf("command %s\n", value)
	}

	for _, value := range blockchain.Spec.ClientArgs {
		log.Printf("ClientArgs %s\n", value)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BlockchainReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&learnv1alpha1.Blockchain{}).
		Complete(r)
}
