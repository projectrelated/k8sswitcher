package kubeconfig

//GetCurrentContext returns "current-context" value
func (k *Kubeconfig) GetCurrentContext() string {
	v := valueOf(k.rootNode, "current-context")
	if v == nil {
		return ""
	}
	return v.Value
}

//UnsetCurrentContext set "current-context" to ""
func (k *Kubeconfig) UnsetCurrentContext() error {
	curCtxValNode := valueOf(k.rootNode, "current-context")
	curCtxValNode.Value = ""
	return nil
}
