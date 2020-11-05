package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var Clientset *kubernetes.Clientset

func InitKubernetesClient(cfgFile string) {
	config, err := clientcmd.BuildConfigFromFlags("", cfgFile)
	if err != nil {
		panic(err)
	}
	Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}
