package kubeconfig

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (k *Kubeconfig) contextsNode() (*yaml.Node, error) {
	contexts := valueOf(k.rootNode, "contexts")
	if contexts == nil {
		return nil, errors.New("\"contexts\" entry is nil")
	} else if contexts.Kind != yaml.SequenceNode {
		return nil, errors.New("\"contexts\" is not in sequence node")
	}
	return contexts, nil
}

func (k *Kubeconfig) contextNode(name string) (*yaml.Node, error) {
	contexts, err := k.contextsNode()
	if err != nil {
		return nil, err

	}

	for _, contextNode := range contexts.Content {
		nameNode := valueOf(contextNode, "name")
		if nameNode.Kind == yaml.ScalarNode && nameNode.Value == name {
			return contextNode, nil
		}

	}
	return nil, errors.Errorf("context with name %q not found", name)
}
