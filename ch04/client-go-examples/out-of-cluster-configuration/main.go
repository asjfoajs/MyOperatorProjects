package main

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"time"
)

func main() {
	homePath := homedir.HomeDir()
	if homePath == "" {
		log.Fatal("home path not found")
	}

	kubeconfig := filepath.Join(homePath, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clinetset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	pods, err := clinetset.CoreV1().Pods("default").
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("There are %d pods in the cluster\n", len(pods.Items))
	for i, pod := range pods.Items {
		log.Printf("%d -> %s/%s", i+1, pod.Namespace, pod.Name)
	}

	<-time.Tick(5 * time.Second)
}
