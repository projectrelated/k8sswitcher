package kubeconfig

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	DefaultLoader Loader = new(StandardKubeconfigLoader)

	Kconfig StandardKubeconfigLoader
)

type ReadWriteResetCloser interface {
	io.ReadWriteCloser

	// Reset truncates the file and seeks to the beginning of the file.
	Reset() error
}

type Loader interface {
	Load() ([]ReadWriteResetCloser, error)
}

type StandardKubeconfigLoader struct{}

type Kubeconfig struct {
	loader Loader

	f        ReadWriteResetCloser
	rootNode *yaml.Node
}

type kubeconfigFile struct{ *os.File }

//HomeDir get the user home directory
func HomeDir() string {
	if v := os.Getenv("XDG_CACHE_HOME"); v != "" {
		return v
	}
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE") // windows
	}
	return home
}

func kubeconfigPath() (string, error) {
	//default path
	home := HomeDir()
	if home == "" {
		return "", errors.New("HOME or USERPROFILE environment variable not set")
	}
	return filepath.Join(home, ".kube", "config"), nil

}

func (*StandardKubeconfigLoader) Load() ([]ReadWriteResetCloser, error) {
	cfgPath, err := kubeconfigPath()
	if err != nil {
		return nil, errors.Wrap(err, "cannot determine kubeconfig path")
	}

	f, err := os.OpenFile(cfgPath, os.O_RDWR, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Wrap(err, "kubeconfig file not found")
		}
		return nil, errors.Wrap(err, "failed to open file")
	}

	return []ReadWriteResetCloser{ReadWriteResetCloser(&kubeconfigFile{f})}, nil
}

func (kf *kubeconfigFile) Reset() error {
	if err := kf.Truncate(0); err != nil {
		return errors.Wrap(err, "failed to truncate file")
	}
	_, err := kf.Seek(0, 0)
	return errors.Wrap(err, "failed to seek in file")
}

func (k *Kubeconfig) WithLoader(l Loader) *Kubeconfig {
	k.loader = l
	return k
}

func (k *Kubeconfig) Parse() error {
	files, err := k.loader.Load()
	if err != nil {
		return errors.Wrap(err, "failed to load")
	}

	f := files[0]

	k.f = f
	var v yaml.Node
	if err := yaml.NewDecoder(f).Decode(&v); err != nil {
		return errors.Wrap(err, "failed to decode")
	}
	k.rootNode = v.Content[0]
	if k.rootNode.Kind != yaml.MappingNode {
		return errors.New("kubeconfig file is not a map document")
	}
	return nil
}

func valueOf(mapNode *yaml.Node, key string) *yaml.Node {
	if mapNode.Kind != yaml.MappingNode {
		return nil
	}
	for i, ch := range mapNode.Content {
		if i%2 == 0 && ch.Kind == yaml.ScalarNode && ch.Value == key {
			return mapNode.Content[i+1]
		}
	}
	return nil
}

func (k *Kubeconfig) ContextNames() []string {
	contexts := valueOf(k.rootNode, "contexts")
	if contexts == nil {
		return nil
	}

	var ctxNames []string
	for _, ctx := range contexts.Content {
		nameVal := valueOf(ctx, "name")
		if nameVal != nil {
			ctxNames = append(ctxNames, nameVal.Value)
		}
	}
	return ctxNames

}

func (k *Kubeconfig) Close() error {
	if k.f == nil {
		return nil
	}
	return k.f.Close()
}

func (k *Kubeconfig) Bytes() ([]byte, error) {
	return yaml.Marshal(k.rootNode)
}

// func main() {

// 	kc := new(Kubeconfig).WithLoader(DefaultLoader)
// 	kc.Parse()

// 	ctxs := kc.ContextNames()
// 	natsort.Sort(ctxs)
// 	for _, c := range ctxs {
// 		fmt.Printf("%s\n", c)
// 	}
// }
