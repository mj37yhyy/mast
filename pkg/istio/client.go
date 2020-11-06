package istio

import (
	k8s "github.com/mj37yhyy/mast/pkg/kubernetes"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	"reflect"
)

var clientset *versionedclient.Clientset

func GetIstioClientset() (*versionedclient.Clientset, error) {
	cs := reflect.ValueOf(clientset)
	if isNil(cs) {
		clientset, err := versionedclient.NewForConfig(k8s.RestConfig)
		return clientset, err
	}
	return clientset, nil
}

func isNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
