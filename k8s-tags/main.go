package main

import (
	"context"
	goError "errors"
	"flag"
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	learnK8SContainer "learn-k8s/container"
	"os"
	"path/filepath"
	"strings"
)


var config *rest.Config
var clientSet *kubernetes.Clientset

func main() {
	useKubeConfig := os.Getenv("USE_KUBECONFIG")
	if len(useKubeConfig) == 0 {
		// default to service account in cluster token
		c, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		config = c
	} else {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		c, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		config = c

	}
	// create the clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	kubeNamespace := os.Getenv("KUBE_NAMESPACE")
	podName := os.Getenv("KUBE_POD_NAME")
	pods, err := clientSet.CoreV1().Pods(kubeNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// Examples for error handling:
	// - Use helper functions like e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	namespace := kubeNamespace
	pod, err := clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", podName, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			podName, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", podName, namespace)
		containers := pod.Spec.Containers
		appContainer, _ := GetAppContainer(containers)

		imageName, shortImage, imageTag, err := learnK8SContainer.SplitImageName(appContainer.Image)
		if err != nil {
			fmt.Errorf("cannot split %s: %s", appContainer.Image, err)
			os.Exit(-1)
		}
		fmt.Printf("cluster_name:%s\n", pod.ClusterName)
		fmt.Printf("kube_namespace:%s\n", pod.Namespace)

		//

		fmt.Println("display_container_name:", fmt.Sprintf("%s_%s", appContainer.Name, pod.Name))
		//fmt.Printf("Container image:%s\n", appContainer.Image)
		fmt.Printf("container_name:%s\n", appContainer.Name)
		fmt.Printf("image_name:%s\n", imageName)
		fmt.Printf("short_image:%s\n", shortImage)
		fmt.Printf("image_tag:%s\n", imageTag)
		fmt.Printf("docker_image:%s\n", fmt.Sprintf("%s:%s", imageName, imageTag))
		fmt.Printf("cloud_provider:alibaba")

		s := "k1=v1; k2=v2; k3=v3"

		entries := strings.Split(s, "; ")

		m := make(map[string]string)
		for _, e := range entries {
			parts := strings.Split(e, "=")
			m[parts[0]] = parts[1]
		}

		fmt.Println(m)
	}

}

func GetAppContainer(containers []v1.Container) (v1.Container, error) {
	for _, c := range containers {
		for _, e := range c.Env {
			//fmt.Printf("%s:%s\n", e.Name, e.Value)
			if e.Name == "DD_AGENT_HOST" {
				return c, nil
			}

		}
	}
	return containers[0], goError.New("could not find App container")
}
