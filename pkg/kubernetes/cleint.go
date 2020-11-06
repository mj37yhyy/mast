package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var (
	Clientset  *kubernetes.Clientset
	RestConfig *restclient.Config
	CfgFile    = filepath.Join(homedir.HomeDir(), ".kube", "config")
)

func init() {
	InitKubernetesClient(CfgFile)
}

func InitKubernetesClient(cfgFile string) {
	if len(cfgFile) == 0 {
		cfgFile = CfgFile
	}
	RestConfig, err := clientcmd.BuildConfigFromFlags("", cfgFile)
	if err != nil {
		panic(err)
	}
	Clientset, err = kubernetes.NewForConfig(RestConfig)
	if err != nil {
		panic(err)
	}
}
