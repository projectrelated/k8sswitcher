package kubeconfig

const (
	defaultNamespace = "default"
)

func (k *Kubeconfig) NamespaceOfContext(contextName string) (string, error) {
	ctx, err := k.contextNode(contextName)
	if err != nil {
		return "", err
	}
	ctxBody := valueOf(ctx, "context")
	if ctxBody == nil {
		return defaultNamespace, nil
	}
	ns := valueOf(ctxBody, "namespace")
	if ns == nil || ns.Value == "" {
		return defaultNamespace, nil
	}
	return ns.Value, nil
}
