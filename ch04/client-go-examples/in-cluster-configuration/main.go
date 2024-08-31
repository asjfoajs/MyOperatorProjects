package main

import (
	"context"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	//1.初始化config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	//2.通过config初始化clientset
	clinetset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	clinetset.AppsV1().Deployments()
	for {
		//3.通过clientset来列出特定命名空间里的所有Pod
		pods, err := clinetset.CoreV1().Pods("default").
			List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("There are %d pods in the cluster\n", len(pods.Items))
		for i, pod := range pods.Items {
			//4.打印信息
			log.Printf("%d -> %s/%s", i+1, pod.Namespace, pod.Name)
		}

		<-time.Tick(5 * time.Second)
	}
}
