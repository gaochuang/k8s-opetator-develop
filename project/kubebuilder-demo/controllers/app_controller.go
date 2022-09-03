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
	"github.com/kubebuilder-demo/controllers/utils"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	v13 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ingressv1beta1 "github.com/kubebuilder-demo/api/v1beta1"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ingress.my.domain,resources=apps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ingress.my.domain,resources=apps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ingress.my.domain,resources=apps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the App object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	app := &ingressv1beta1.App{}
	//从缓存中获取APP
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	//根据App配置进行处理
	//1.Deployment处理
	deployment := utils.NewDeployment(app)
	//设置所属对象的绑定，当owner被删除时，与其绑定的资源也会被删除，App被删除时，deployment也会被删除
	if err := controllerutil.SetControllerReference(app, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	//2. 查找同名的deployment
	d := &v1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, d)
	if err != nil {
		if errors.IsNotFound(err) {
			err := r.Create(ctx, deployment)
			if err != nil {
				logger.Error(err, "create deployment failed")
				return ctrl.Result{}, err
			}
		} else {
			err := r.Update(ctx, deployment)
			if err != nil {
				logger.Error(err, "Update deployment failed")
				return ctrl.Result{}, err
			}
		}
	}

	//3.service处理
	//设置所属对象的绑定，当owner被删除时，与其绑定的资源也会被删除，App被删除时，service也会被删除
	service := utils.NewService(app)
	if err := controllerutil.SetControllerReference(app, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	//4.查找service
	s := &v12.Service{}
	err = r.Get(ctx, req.NamespacedName, s)
	if err != nil {
		if errors.IsNotFound(err) && app.Spec.EnableService {
			err := r.Create(ctx, service)
			if err != nil {
				logger.Error(err, "Create service failed")
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
			err := r.Delete(ctx, s)
			if err != nil {
				logger.Error(err, "Delete service failed")
				return ctrl.Result{}, err
			}
		}
	}

	//创建Ingress的前提是必须创建service
	//Todo 使用admission设置默认值,默认为false
	if !app.Spec.EnableService {
		return ctrl.Result{}, nil
	}

	//5.Ingress处理
	////设置所属对象的绑定，当owner被删除时，与其绑定的资源也会被删除，App被删除时，service也会被删除
	ingress := utils.NewIngress(app)
	err = controllerutil.SetControllerReference(app, ingress, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	in := &v13.Ingress{}
	err = r.Get(ctx, req.NamespacedName, in)
	if err != nil {
		if errors.IsNotFound(err) && app.Spec.EnableIngress {
			err := r.Create(ctx, ingress)
			if err != nil {
				logger.Error(err, "Create ingress failed")
				return ctrl.Result{}, err
			}
		}
		if !errors.IsNotFound(err) && app.Spec.EnableIngress {
			return ctrl.Result{}, err
		}
	} else {
		if app.Spec.EnableIngress {
			logger.Info("skip update")
		} else {
			if err := r.Delete(ctx, in); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
//Manager的client传给Controller,并且调用SetupWithManager方法传入给Manager进行controller的初始化
func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ingressv1beta1.App{}).
		Owns(&v1.Deployment{}).
		Owns(&v12.Service{}).
		Owns(&v13.Ingress{}).
		Complete(r)
}
