package kubernetes

import (
	"fmt"
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
	if err := InitKubernetesClient(CfgFile); err != nil {
		fmt.Printf("Please copy kubernetes config file to %s or call \"mast config <your_path>/config\". \n", CfgFile)
	}
}

func InitKubernetesClient(cfgFile string) error {
	if len(cfgFile) == 0 {
		cfgFile = CfgFile
	}
	RestConfig, err := clientcmd.BuildConfigFromFlags("", cfgFile)
	if err != nil {
		return err
	}
	Clientset, err = kubernetes.NewForConfig(RestConfig)
	if err != nil {
		return err
	}
	return nil
}
