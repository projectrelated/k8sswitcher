package kubeconfig

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	//"github.com/projectrelated/k8sswitcher/kubeconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func Run(stdout, stderr io.Writer) error {
	kc := new(Kubeconfig).WithLoader(DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return errors.Wrap(err, "kubeconfig error")
	}

	ctx := kc.GetCurrentContext()
	if ctx == "" {
		return errors.New("current-context is not set")
	}
	curNs, err := kc.NamespaceOfContext(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot read current namespace")
	}

	ns, err := QueryNamespaces(kc)
	if err != nil {
		return errors.Wrap(err, "could not list namespaces (is the cluster accessible?)")
	}

	for _, c := range ns {
		s := c
		if c == curNs {
			fmt.Println("CurrentNamespace", s) //Will need to go to systray as a tick mark of the current context
		}
	}
	return nil
}

func QueryNamespaces(kc *Kubeconfig) ([]string, error) {
	if os.Getenv("_MOCK_NAMESPACES") != "" {
		return []string{"ns1", "ns2"}, nil
	}

	clientset, err := newKubernetesClientSet(kc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize k8s REST client")
	}

	var out []string
	var next string
	for {
		list, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{
			Limit:    500,
			Continue: next,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list namespaces from k8s API")
		}
		next = list.Continue
		for _, it := range list.Items {
			out = append(out, it.Name)
		}
		if next == "" {
			break
		}
	}
	return out, nil
}

func newKubernetesClientSet(kc *Kubeconfig) (*kubernetes.Clientset, error) {
	b, err := kc.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert in-memory kubeconfig to yaml")
	}
	cfg, err := clientcmd.RESTConfigFromKubeConfig(b)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize config")
	}
	return kubernetes.NewForConfig(cfg)
}
