package main

import "k8s.io/client-go/tools/clientcmd"

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedConfigPathFlag)
	if err != nil {
		panic(err.Error())
	}

}
