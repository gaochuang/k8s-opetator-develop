package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	fmt.Println(clientcmd.RecommendedHomeFile)

	//config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	config.APIPath = "api"                           //请求的HTTP路径
	config.GroupVersion = &corev1.SchemeGroupVersion //请求资源的版本
	config.NegotiatedSerializer = scheme.Codecs      //数据编解码器

	//client
	restclient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	//get data
	result := &corev1.PodList{}

	err = restclient.Get().
		Namespace("default").
		Resource("pods").
		Do(context.TODO()).
		Into(result)

	if err != nil {
		panic(err)
	}
	for _, d := range result.Items {
		fmt.Printf("NAMESPACE: %v \t NAME:%v \t Status:%+v\n", d.Namespace, d.Name, d.Status.Phase)
	}

}
