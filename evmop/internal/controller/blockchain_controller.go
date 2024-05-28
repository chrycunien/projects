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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	apiResource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *BlockchainReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.SetPrefix("BlockchainReconciler")

	blockchain := &learnv1alpha1.Blockchain{}
	err := r.Get(ctx, req.NamespacedName, blockchain)
	if err != nil {
		return ctrl.Result{}, nil
	}

	/*
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
	*/

	// Check if the statefulset already exists, if not create a new one
	foundSts := &appsv1.StatefulSet{}
	err = r.Get(
		context.Background(),
		types.NamespacedName{
			Name:      blockchain.Name,
			Namespace: blockchain.Namespace,
		},
		foundSts,
	)
	if err != nil && errors.IsNotFound(err) {
		// Create a new StatefulSet
		sts := r.ReconcileStatefulSet(blockchain)
		err = r.Client.Create(context.Background(), sts)
		if err != nil {
			log.Println("Failed to create new StatefulSet", err, "NameSpace", sts.Namespace, "Name", sts.Name)
			return ctrl.Result{}, nil
		}
		// Statefulset created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Println("Failed to get StatefulSet", err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *BlockchainReconciler) ReconcileStatefulSet(b *learnv1alpha1.Blockchain) *appsv1.StatefulSet {
	log.Println("Creating a new StatefulSet")

	// Make sure to run at least 1 replicas
	if b.Spec.Replicas == nil {
		b.Spec.Replicas = pointer.Int32(1)
	}

	// Provisioning a PVC to store this statefulset's data
	pvc := v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "data",
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			StorageClassName: pointer.String("standard"),
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: apiResource.MustParse("1Gi"),
				},
			},
		},
	}

	if b.Spec.Cpu == "" {
		b.Spec.Cpu = "500m"
	}
	if b.Spec.Memory == "" {
		b.Spec.Memory = "1Gi"
	}

	// Specifying resources for the main container
	reqs := &v1.ResourceRequirements{
		Limits: v1.ResourceList{
			"cpu":    apiResource.MustParse(b.Spec.Cpu),
			"memory": apiResource.MustParse(b.Spec.Memory),
		},
		Requests: v1.ResourceList{
			"cpu":    apiResource.MustParse(b.Spec.Cpu),
			"memory": apiResource.MustParse(b.Spec.Memory),
		},
	}

	if b.Spec.ApiPort == 0 {
		b.Spec.ApiPort = 8545
	}

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.Name,
			Namespace: b.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: b.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: b.ObjectMeta.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: b.ObjectMeta.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image:           b.Spec.Image,
							ImagePullPolicy: "Always",
							Name:            "app",
							Command:         b.Spec.Command,
							Args:            b.Spec.ClientArgs,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 30303,
									Name:          "p2p",
									Protocol:      "TCP",
								},
								{
									ContainerPort: b.Spec.ApiPort,
									Name:          "api",
									Protocol:      "TCP",
								},
							},
							Resources: *reqs,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/data",
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				pvc,
			},
		},
	}

	// Set Learn instance as the owner and controller
	controllerutil.SetControllerReference(b, sts, r.Scheme)

	return sts
}

// SetupWithManager sets up the controller with the Manager.
func (r *BlockchainReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&learnv1alpha1.Blockchain{}).
		Complete(r)
}
