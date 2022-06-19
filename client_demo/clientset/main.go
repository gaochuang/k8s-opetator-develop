package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
)

const (
	NAMESPACE       = "test-namespace"
	SERVICE_NAME    = "client-test-service"
	DEPLOYMENT_NAME = "client-test-deployment"
)

func main() {

	restconfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err.Error())
	}

	operate := flag.String("operate", "create", "opetate type: create or clean")

	flag.Parse()

	//client是多个客户端的集合
	clientSet, err := kubernetes.NewForConfig(restconfig)
	if err != nil {
		panic(err.Error())
	}

	switch *operate {
	case "create":
		createSource(clientSet)
	case "clean":
		cleanResource(clientSet)
	default:
		fmt.Println("don't support this operate, please use create or clean")
	}
}

//create source
func createSource(clientSet *kubernetes.Clientset) {
	//create namespace
	createNamespace(clientSet)
	//create svc
	createService(clientSet)
	//createDeployment
	createDeployment(clientSet)
}

func createNamespace(clientSet *kubernetes.Clientset) {

	//1.创建一个namespace的client  client-go/kubernetes/typed/core/v1/core_client.go
	namespaceClient := clientSet.CoreV1().Namespaces()
	//2.初始化namespac的数据结构，kubernetes/staging/src/k8s.io/api/core/v1/types.go
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: NAMESPACE,
		},
	}
	//3.实例化namespace, client-go/kubernetes/typed/core/v1/namespace.go
	ret, err := namespaceClient.Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Create namespace is %s \n", ret.GetName())
}

func createService(clientSet *kubernetes.Clientset) {

	//1.创建service client client-go/kubernetes/typed/core/v1/core_client.go
	serviceClient := clientSet.CoreV1().Services(NAMESPACE)

	//2.初始化service, kubernetes/staging/src/k8s.io/api/core/v1/types.go
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: SERVICE_NAME,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:     "http",
				Port:     8080,
				NodePort: 30080,
			},
			},
			Selector: map[string]string{
				"testapp": "tomcat",
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	//3.实例化service, client-go/kubernetes/typed/core/v1/service.go
	ret, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	log.Printf("Create service %s \n", ret.GetName())
}

func createDeployment(clientSet *kubernetes.Clientset) {

	//1.创建Deployment控制器的client(api属于apps), client-go/kubernetes/typed/apps/v1/apps_client.go
	deploymentClient := clientSet.AppsV1().Deployments(NAMESPACE)

	//2.初始化deployment，code/kubernetes/staging/src/k8s.io/api/apps/v1/types.go
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: DEPLOYMENT_NAME,
		},

		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"testapp": "tomcat",
				},
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"testapp": "tomcat",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "tomcat",
							Image:           "tomcat",
							ImagePullPolicy: "IfNotPresent",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolSCTP,
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}

	//3.实例化Deployment
	ret, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	log.Printf("create deployment %s\n", ret.GetName())
}

func cleanResource(clientSet *kubernetes.Clientset) {

	emptyDeleteOptions := metav1.DeleteOptions{}

	//delete service
	if err := clientSet.CoreV1().Services(NAMESPACE).Delete(context.TODO(), SERVICE_NAME, emptyDeleteOptions); err != nil {
		log.Printf("delete service failed %s\n", SERVICE_NAME)
	}
	//delete deplotment
	if err := clientSet.AppsV1().Deployments(NAMESPACE).Delete(context.TODO(), DEPLOYMENT_NAME, emptyDeleteOptions); err != nil {
		log.Printf("delete deployment failed %s\n", DEPLOYMENT_NAME)
	}

	//delete namespace
	if err := clientSet.CoreV1().Namespaces().Delete(context.TODO(), NAMESPACE, emptyDeleteOptions); err != nil {
		log.Printf("delete namespace failed %s\n", NAMESPACE)
	}

}
