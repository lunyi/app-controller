/*
Copyright 2023.

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

	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appv1 "app-controller/api/v1"
	"app-controller/internal/controller/utils"
)

// LobbyReconciler reconciles a Lobby object
type LobbyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.srvptt2.online,resources=lobbies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=service,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.srvptt2.online,resources=lobbies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.srvptt2.online,resources=lobbies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Lobby object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *LobbyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	logger := log.FromContext(ctx)

	logger.Info("Start LobbyController")

	app := &appv1.Lobby{}
	//从缓存中获取app
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	res, err := createDeployment(app, logger, r, ctx, &req)

	if err != nil {
		return ctrl.Result{}, err
	}

	res, err = createService(app, logger, r, ctx)

	if err != nil {
		return ctrl.Result{}, err
	}

	res, err = createIngress(app, logger, r, ctx)

	if err != nil {
		return ctrl.Result{}, err
	}

	return res, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LobbyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.Lobby{}).
		Owns(&v1.Deployment{}).
		Owns(&netv1.Ingress{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

func createDeployment(
	app *appv1.Lobby,
	logger logr.Logger,
	r *LobbyReconciler,
	ctx context.Context,
	req *reconcile.Request) (ctrl.Result, error) {

	deployment := utils.NewDeployment(app, logger)
	if err := controllerutil.SetControllerReference(app, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	//查找同名deployment
	d := &v1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, d); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, deployment); err != nil {
				logger.Error(err, "create deploy failed")
				return ctrl.Result{}, err
			}
		}
	} else {
		if app.Spec.Replicas != int(*deployment.Spec.Replicas) &&
			app.Spec.Image != deployment.Spec.Template.Spec.Containers[0].Image {

			if err := r.Update(ctx, deployment); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}

func createService(
	app *appv1.Lobby,
	logger logr.Logger,
	r *LobbyReconciler,
	ctx context.Context) (ctrl.Result, error) {

	service := utils.NewService(app, logger)
	if err := controllerutil.SetControllerReference(app, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	s := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{Name: app.Name, Namespace: app.Namespace}, s); err != nil {
		if errors.IsNotFound(err) && app.Spec.EnableService {
			if err := r.Create(ctx, service); err != nil {
				logger.Error(err, "create service failed")
				return ctrl.Result{}, err
			}
		}
		if !errors.IsNotFound(err) && app.Spec.EnableService {
			return ctrl.Result{}, err
		}
	} else {
		if app.Spec.EnableService {
			logger.Info("skip update service")
		} else {
			if err := r.Delete(ctx, s); err != nil {
				return ctrl.Result{}, err
			}

		}
	}
	return ctrl.Result{}, nil
}

func createIngress(
	app *appv1.Lobby,
	logger logr.Logger,
	r *LobbyReconciler, // Assuming r is an instance of LobbyReconciler
	ctx context.Context,
) (ctrl.Result, error) { // Fixing return type to ctrl.Result and error
	ingress := utils.NewIngress(app, logger)
	if err := controllerutil.SetControllerReference(app, ingress, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	i := &netv1.Ingress{}
	if err := r.Get(ctx, types.NamespacedName{Name: app.Name, Namespace: app.Namespace}, i); err != nil {
		if errors.IsNotFound(err) && app.Spec.EnableIngress {
			if err := r.Create(ctx, ingress); err != nil {
				logger.Error(err, "create ingress failed")
				return ctrl.Result{}, err
			}
		}
		if !errors.IsNotFound(err) && app.Spec.EnableIngress {
			return ctrl.Result{}, err
		}
	} else {
		if app.Spec.EnableIngress {
			logger.Info("skip update ingress")
		} else {
			if err := r.Delete(ctx, i); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}
