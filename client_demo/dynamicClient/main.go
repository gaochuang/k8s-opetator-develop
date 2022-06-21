package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err.Error())
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//resource schema.GroupVersionResource
	//dynamicClient的唯一关联方法所需的入参
	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	//使用dynamicClient方法查询列表方法，查询指定namespace下所有的Pod
	//返回的数据结构类型是UnstructuredList
	unstructObj, err := dynamicClient.Resource(gvr).
		Namespace("default").
		List(context.TODO(), metav1.ListOptions{Limit: 100})

	if err != nil {
		panic(err.Error())
	}

	//实例化一个PodList用于存放unstructObj转换后的结果
	PodList := &corev1.PodList{}

	//进行转换
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), PodList); err != nil {
		panic(err.Error())
	}

	//每个Pod都打印出namespace、status.Phase\name三个字段
	for _, d := range PodList.Items {
		fmt.Printf("%v\t %v\t %v\n",
			d.Namespace,
			d.Status.Phase,
			d.Name)
	}
}
