package main

import (
	"pkg"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//1. create config
	//从集群外部查找配置文件
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		//内部创建，通过serviceAccount挂载的文件创建
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	//2. create clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//3. 创建factory
	factoryInformers := informers.NewSharedInformerFactory(clientSet, 0)

	//4. 创建service informer
	serviceInformer := factoryInformers.Core().V1().Services()

	//5.创建ingree informer
	ingressInformer := factoryInformers.Networking().V1().Ingresses()

	//6.创建controller
	controller := pkg.NewController(clientSet, serviceInformer, ingressInformer)

	stopCh := make(chan struct{})

	//7.启动informer
	factoryInformers.Start(stopCh)

	//8.启动同步
	factoryInformers.WaitForCacheSync(stopCh)

	controller.Run(stopCh)

}
