module github.com/projectrelated/k8sswitcher

go 1.13

replace github.com/projectrelated/k8sswitcher/k8sclient => /k8sclient

require (
	facette.io/natsort v0.0.0-20181210072756-2cd4dd1e2dcb
	github.com/ahmetb/kubectx v0.9.1 // indirect
	github.com/getlantern/systray v0.0.0-20200324212034-d3ab4fd25d99
	github.com/goccy/go-yaml v1.8.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79 // indirect
	golang.org/x/net v0.0.0-20200501053045-e0ff5e5a1de5 // indirect
	gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	k8s.io/api v0.17.5
	k8s.io/apimachinery v0.17.5
	k8s.io/client-go v0.17.5
	sigs.k8s.io/yaml v1.2.0
)
