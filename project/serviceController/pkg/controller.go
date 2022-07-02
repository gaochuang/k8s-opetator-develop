package pkg

import (
	"context"
	"fmt"
	v12 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	v13 "k8s.io/client-go/informers/core/v1"
	v14 "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	coreV1 "k8s.io/client-go/listers/core/v1"
	V1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"time"
)

const (
	worker   = 5
	maxRetry = 10
)

type controller struct {
	client        kubernetes.Interface
	ingressLister V1.IngressLister
	serviceLister coreV1.ServiceLister
	//限速队列
	queue workqueue.RateLimitingInterface
}

func (c *controller) addService(obj interface{}) {
	c.enqueue(obj)
}

func (c *controller) updateService(oldobj interface{}, newobj interface{}) {
	//update检查annotation是否发生变化，如果没有更新不用处理
	if reflect.DeepEqual(oldobj, newobj) {
		return
	}
	c.enqueue(newobj)
}

//统一接口，将事件放入限速队列
func (c *controller) enqueue(obj interface{}) {
	//查找obj的key
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}
	c.queue.Add(key)
}

func (c *controller) deleteIngress(obj interface{}) {
	ingress := obj.(*v12.Ingress)
	//获取ingress对应的service
	ownerReference := v1.GetControllerOf(ingress)

	//如果service不存在，就不用处理
	if ownerReference == nil {
		return
	}

	//如果不是service
	if ownerReference.Kind != "service" {
		return
	}

	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
}

func (c *controller) Run(ch chan struct{}) {
	fmt.Println("Controller start")
	for i := 0; i < worker; i++ {
		go wait.Until(c.worker, time.Minute, ch)
	}
}

func (c *controller) worker() {
	for c.processNextItem() {

	}
}

func (c *controller) processNextItem() bool {

	//获取key
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	key := item.(string)
	//处理完将事件从对象删除
	defer c.queue.Done(item)

	err := c.syncService(key)
	if err != nil {
		c.handlerError(key, err)
	}

	return true
}

func (c *controller) syncService(item string) error {

	//分离service的namespace、name
	namespaceKey, name, err := cache.SplitMetaNamespaceKey(item)
	if err != nil {
		return err
	}

	//判断service是否存在
	service, err := c.serviceLister.Services(namespaceKey).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	//判断annotation ingress/http是否存在
	_, ok := service.GetAnnotations()["ingress/http"]

	//判断ingress是否存在
	_, err = c.ingressLister.Ingresses(namespaceKey).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	//如果service存在，但是ingress不存在，需要创建ingress
	if ok && errors.IsNotFound(err) {
		//创建ingress
		ig := c.createIngress(namespaceKey, name)
		_, err := c.client.NetworkingV1().Ingresses(namespaceKey).Create(context.TODO(), ig, v1.CreateOptions{})
		if err != nil {
			return err
		}
	} else if !ok && !errors.IsNotFound(err) {
		//如果service不存在，但是ingress存在需要删除ingress
		err = c.client.NetworkingV1().Ingresses(namespaceKey).Delete(context.TODO(), name, v1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil

}

func (c *controller) createIngress(namespace string, name string) *v12.Ingress {

	ingress := v12.Ingress{}
	ingress.Name = name
	ingress.Namespace = namespace
	pathType := v12.PathTypePrefix
	ingress.Spec = v12.IngressSpec{
		Rules: []v12.IngressRule{
			{
				Host: "example.com",
				IngressRuleValue: v12.IngressRuleValue{
					HTTP: &v12.HTTPIngressRuleValue{
						Paths: []v12.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: v12.IngressBackend{
									Service: &v12.IngressServiceBackend{
										Name: name,
										Port: v12.ServiceBackendPort{
											Number: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return &ingress
}

func (c *controller) handlerError(key string, err error) {
	//如果出现错误，再次将出错的key加入到限速,限制只能够重试10次
	if c.queue.NumRequeues(key) <= maxRetry {
		c.queue.AddRateLimited(key)
		return
	}
	runtime.HandleError(err)
	//不记录重试次数
	c.queue.Forget(key)

}

func NewController(client kubernetes.Interface, serviceInformer v13.ServiceInformer, ingressInformer v14.IngressInformer) controller {

	c := controller{
		client:        client,
		serviceLister: serviceInformer.Lister(),
		ingressLister: ingressInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}

	serviceInformer.Informer().AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})

	ingressInformer.Informer().AddEventHandler(&cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})

	return c
}
