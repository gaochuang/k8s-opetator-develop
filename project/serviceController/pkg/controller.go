package pkg

import (
	V13 "k8s.io/client-go/informers/core/v1"
	V12 "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	coreV1 "k8s.io/client-go/listers/core/v1"
	V1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
)

type controller struct {
	client        kubernetes.Interface
	ingressLister V1.IngressLister
	serviceLister coreV1.ServiceLister
}

func (c *controller) addService(obj interface{}) {

}

func (c *controller) updateService(old interface{}, new interface{}) {

}

func (c *controller) deleteIngress(obj interface{}) {

}

func (c *controller) Run() {

}

func NewController(client kubernetes.Interface, serviceInformer V13.ServiceInformer, ingressInformer V12.IngressInformer) controller {

	c := controller{
		client:        client,
		serviceLister: serviceInformer.Lister(),
		ingressLister: ingressInformer.Lister(),
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
