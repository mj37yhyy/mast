package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var (
	Clientset *kubernetes.Clientset
	CfgFile   = filepath.Join(homedir.HomeDir(), ".kube", "config")
)

func init() {
	InitKubernetesClient(CfgFile)
}

func InitKubernetesClient(cfgFile string) {
	if len(cfgFile) == 0 {
		cfgFile = CfgFile
	}
	config, err := clientcmd.BuildConfigFromFlags("", cfgFile)
	if err != nil {
		panic(err)
	}
	Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}
