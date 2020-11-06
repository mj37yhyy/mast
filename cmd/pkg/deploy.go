package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/mj37yhyy/mast/pkg/istio"
	k8s "github.com/mj37yhyy/mast/pkg/kubernetes"
	"github.com/spf13/cobra"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"strings"
)

var (
	namespace      string
	serviceName    string
	serviceVersion string
	image          string
	ports          []string
	Ports          [][]int32
	//nodePort       int32
	replicas  int32
	deployCmd = &cobra.Command{
		Use:               "deploy",
		Short:             "Deploy the service to kubernetes and use istio to publish in grayscale.",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(image) != 0 && len(serviceName) != 0 && len(serviceVersion) != 0 {
				if !printErr(cmd, initPorts) ||
					!printErr(cmd, createDeployment) ||
					!printErr(cmd, createOrUpdateService) ||
					!printErr(cmd, createOrUpdateDestinationRule) {
					return
				}
			} else {
				if !printErr(cmd, cmd.Help) {
					return
				}
			}
		},
	}
)

func init() {
	deployCmd.PersistentFlags().StringVar(&image, "image", "", "docker image name for kubernetes.")
	deployCmd.PersistentFlags().StringArrayVar(&ports, "ports", []string{"8080:8080"},
		"Define multiple port. The format is <container port>:<service port>. "+
			"example: -ports 8080:8080 -ports 443:443. default 8080:8080.")
	//deployCmd.PersistentFlags().Int32Var(&containerPort, "targetPort", 8080, "container port. default 8080")
	//deployCmd.PersistentFlags().Int32Var(&servicePort, "port", 8080, "service port. default 8080")
	//deployCmd.PersistentFlags().Int32Var(&nodePort, "nodePort", 0, "container port")
	deployCmd.PersistentFlags().Int32Var(&replicas, "replicas", 1, "replicas for pod. default 1")
	deployCmd.PersistentFlags().StringVar(&namespace, "namespace", apiv1.NamespaceDefault,
		"namespace for kubernetes. default default")
	deployCmd.PersistentFlags().StringVar(&serviceName, "name", "", "deployment name for kubernetes")
	deployCmd.PersistentFlags().StringVar(&serviceVersion, "version", "", "deployment version for istio")
}

func initPorts() error {
	for i := range ports {
		portStringArray := strings.Split(ports[i], ":")
		if len(portStringArray) != 2 {
			return errors.New("ports format error.")
		}
		var portInt32Array []int32

		p, err := strconv.ParseInt(portStringArray[0], 10, 32)
		if err != nil {
			panic(err)
		}
		portInt32Array = append(portInt32Array, int32(p))

		p, err = strconv.ParseInt(portStringArray[1], 10, 32)
		if err != nil {
			panic(err)
		}
		portInt32Array = append(portInt32Array, int32(p))

		Ports = append(Ports, portInt32Array)
	}
	return nil
}

func createDeployment() error {
	deploymentsClient := k8s.Clientset.AppsV1().Deployments(namespace)
	var dPorts []apiv1.ContainerPort
	for i := range Ports {
		dPorts = append(dPorts, apiv1.ContainerPort{
			Protocol:      apiv1.ProtocolTCP,
			ContainerPort: Ports[i][0],
		})
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     serviceName,
					"version": serviceVersion,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     serviceName,
						"version": serviceVersion,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							//Name:  "web",
							Image: image,
							Ports: dPorts,
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return nil
}

func createOrUpdateService() error {
	svcClient := k8s.Clientset.CoreV1().Services(namespace)

	var sPorts []apiv1.ServicePort
	for i := range Ports {
		sPorts = append(sPorts, apiv1.ServicePort{
			TargetPort: intstr.IntOrString{
				IntVal: Ports[i][0],
			},
			Port: Ports[i][1],
		})
	}

	svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": serviceName,
			},
			Ports: sPorts,
		},
	}

	svc, err := svcClient.Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		svc, err = svcClient.Update(context.TODO(), svc, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func createOrUpdateDestinationRule() error {
	istioClientset, err := istio.GetIstioClientset()
	if err != nil {
		return err
	}
	dri := istioClientset.NetworkingV1beta1().DestinationRules(namespace)

	drName := serviceName + "." + namespace
	dr := &v1beta1.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: networkingv1beta1.DestinationRule{
			Host: drName,
			Subsets: []*networkingv1beta1.Subset{
				{
					Name: serviceVersion,
					Labels: map[string]string{
						version: serviceVersion,
					},
				},
			},
		},
	}

	dr, err = dri.Create(context.TODO(), dr, metav1.CreateOptions{})
	if err != nil {
		pdr := fmt.Sprintf(`spec:
  subsets:
  - name: %s
    labels:
      version: %s`, serviceVersion, serviceVersion)
		dr, err = dri.Patch(context.TODO(), drName, types.MergePatchType, []byte(pdr), metav1.PatchOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }

func printErr(cmd *cobra.Command, _func func() error) bool {
	if err := _func(); err != nil {
		cmd.PrintErr(err)
		return false
	}
	return true
}
