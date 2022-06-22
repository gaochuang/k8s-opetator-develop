package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err.Error())
	}

	discoverClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	APIGroup, APIResourceList, err := discoverClient.ServerGroupsAndResources()
	if err != nil {
		panic(err.Error())
	}

	//打印Group信息
	fmt.Printf("APIGroup: \n\n%v\n\n", APIGroup)

	//
	for _, sigAPIResouceList := range APIResourceList {

		//GroupVersion是个字符串 "apps/v1"
		groupVersionStr := sigAPIResouceList.GroupVersion

		//将字符串转换成数据结构
		gv, err := schema.ParseGroupVersion(groupVersionStr)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("########################################################")
		fmt.Printf("GV string [%v]\n Gv struct [%#v]\nresources: \n\n", groupVersionStr, gv)
		for _, singelAPIreource := range sigAPIResouceList.APIResources {
			fmt.Printf("%v\n", singelAPIreource.Name)
		}

	}

}
