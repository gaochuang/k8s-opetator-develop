package main

import (
	"context"
	v1 "create-crd-template/pkg/apis/mycontoller.temp/v1"
	"fmt"
	"log"

	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalln(err)
	}
	config.APIPath = "/apis/"
	config.GroupVersion = &v1.GroupVersion
	//指定解码工具
	config.NegotiatedSerializer = v1.Codes.WithoutConversion()

	client, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatalln(err)
	}

	foo := v1.Foo{}
	err = client.Get().Namespace("default").Resource("foos").Name("test-crd").Do(context.TODO()).Into(&foo)
	if err != nil {
		log.Fatalln(err)
	}
	newObj := foo.DeepCopy()
	newObj.Spec.Name = "test2"
	fmt.Println(foo.Spec.Name)
	fmt.Println(foo.Spec.Replicas)
	fmt.Println(newObj.Spec.Name)
}
