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
	learnK8SFile "learn-k8s/file"
	"os"
	"path/filepath"
	"strings"
)

const fileName = "output-datadog.conf"
const outputFolder = "dd-output-config"

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
		var tags []string
		imageName, shortImage, imageTag, err := learnK8SContainer.SplitImageName(appContainer.Image)
		if err != nil {
			fmt.Errorf("cannot split %s: %s", appContainer.Image, err)
			os.Exit(-1)
		}


		tags = append(tags, "cluster_name:" + pod.ClusterName)
		tags = append(tags, "kube_namespace:" + pod.Namespace)
		tags = append(tags, "pod_name:" + pod.Name)

		//fmt.Printf("cluster_name:%s\n", pod.ClusterName)
		//fmt.Printf("kube_namespace:%s\n", pod.Namespace)
		//fmt.Println("display_container_name:", fmt.Sprintf("%s_%s", appContainer.Name, pod.Name))
		////fmt.Printf("Container image:%s\n", appContainer.Image)
		//fmt.Printf("container_name:%s\n", appContainer.Name)
		//fmt.Printf("image_name:%s\n", imageName)
		//fmt.Printf("short_image:%s\n", shortImage)
		//fmt.Printf("image_tag:%s\n", imageTag)
		//fmt.Printf("docker_image:%s\n", fmt.Sprintf("%s:%s", imageName, imageTag))
		//fmt.Printf("cloud_provider:alibaba")

		tags = append(tags, "display_container_name:" + fmt.Sprintf("%s_%s", appContainer.Name, pod.Name))
		tags = append(tags, "kube_container_name:" + appContainer.Name)
		tags = append(tags, "container_id:" + getContainerIdByName(appContainer.Name, pod))
		tags = append(tags, "image_name:" + imageName)
		tags = append(tags, "short_image:" + shortImage)
		tags = append(tags, "image_tag:" + imageTag)
		tags = append(tags, "docker_image:" + fmt.Sprintf("%s:%s", imageName, imageTag))
		tags = append(tags, "cloud_provider:" + "alibaba")
		tags = append(tags, "host:" + "datadog-agent-78b568f767-z2z2x")



		//learnK8SFile.ReadFile(strings.Join(tags, ","), fileName)
		fmt.Printf("K8S Tags:" + strings.Join(tags, ",")+"\n")
		if learnK8SFile.CheckFileExists(fileName) && learnK8SFile.CheckFileExists("parsers.conf") && learnK8SFile.CheckFileExists("fluent-bit.conf") {
			//learnK8SFile.AppendFile(strings.Join(tags, ","), fileName)
			learnK8SFile.CreateDirectoryIfNotExists(outputFolder)
			learnK8SFile.WriteFile(outputFolder+"/"+fileName, learnK8SFile.ReadFile(strings.Join(tags, ","), fileName))
			learnK8SFile.CreateACopyOfFile("parsers.conf", outputFolder+"/"+"parsers.conf")
			learnK8SFile.CreateACopyOfFile("fluent-bit.conf", outputFolder+"/"+"fluent-bit.conf")

		} else {
			learnK8SFile.GetCurrentDirectory()
			learnK8SFile.PrintFilesInDirectory()
		}

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

func getContainerIdByName(containerName string, pod *v1.Pod) string {
	for _, c := range pod.Status.ContainerStatuses {
		if c.Name == containerName {
			return c.ContainerID
		}
	}
	return ""
}