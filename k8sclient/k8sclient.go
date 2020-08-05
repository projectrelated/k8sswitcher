package k8sclient

import (
	"flag"
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/mitchellh/go-homedir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//GetNamespace ... get the names space in k8s
func GetNamespace() (namespacelist []string) {

	var label, field string
	flag.StringVar(&label, "l", "", "Label selector")
	flag.StringVar(&field, "f", "", "Field selector")

	homedir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	kubeconfig := filepath.Join(
		homedir, ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	api := clientset.CoreV1()
	listOptions := metav1.ListOptions{
		LabelSelector: label,
		FieldSelector: field,
	}
	namespaces, err := api.Namespaces().List(listOptions)
	if err != nil {
		log.Fatal(err)
	}
	var namespaceslice []string
	for _, names := range namespaces.Items {
		namespaceslice = append(namespaceslice, names.ObjectMeta.Name)
		fmt.Println(names.ObjectMeta.Name)
	}
	fmt.Println("namespaceslice", namespaceslice)
	return namespaceslice
}
