package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"facette.io/natsort"
	"github.com/getlantern/systray"
	"github.com/projectrelated/k8sswitcher/kubeconfig"
	log "github.com/sirupsen/logrus"
	//	k8sclient "github.com/projectrelated/k8sswitcher/k8sclient"
)

type ClusterConfig struct {
	Clusters []struct {
		Name    string `yaml:"name"`
		Cluster struct {
			Server string `yaml:"server"`
		} `yaml:"cluster"`
	} `yaml:"clusters"`
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}

var file *os.File

var clusterMap map[string][]string

func main() {

	//Set Log level to debug using LogLevel Environment Variable. ex: DebugLevel
	envLogLevel := os.Getenv("LogLevel")

	logFile, err := os.OpenFile("k8sswitcher.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetLevel(log.InfoLevel)
	if envLogLevel == "DebugLevel" {
		log.SetLevel(log.DebugLevel)
		log.Info("Setting log level to Debug")
	}

	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	kc.Parse()

	CurrentContext := kc.GetCurrentContext()
	CurrentNamespace, err := kc.NamespaceOfContext(CurrentContext)
	if err != nil {
		fmt.Println("failed")
	}
	Namespaces, err := kubeconfig.QueryNamespaces(kc)

	fmt.Println(CurrentContext)
	fmt.Println(CurrentNamespace)
	fmt.Println(Namespaces)

	ctxs := kc.ContextNames()
	natsort.Sort(ctxs)
	for _, c := range ctxs {
		fmt.Printf("%s\n", c)
	}

	clusterMap = make(map[string][]string)

	for _, contextName := range ctxs {
		clusterMap[contextName] = []string{"test-namespace"}
	}

	namespaceMap := []string{}
	for _, namespace := range Namespaces {

		namespaceMap = append(namespaceMap, namespace)
	}

	clusterMap["docker-desktop"] = namespaceMap

	//clusterMap["test-cluster-1"] = []string{"clus1-namespace1", "clus1-namespace2", "clus1-namespace3", "clus1-namespace4"}
	//clusterMap["test-cluster-2"] = []string{"clus2-namespace1", "clus2-namespace2", "clus2-namespace3", "clus2-namespace4"}

	onExit := func() {
		fmt.Println("Starting onExit")
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
		fmt.Println("Finished onExit")
	}
	// Should be called at the very beginning of main().
	systray.RunWithAppWindow("Lantern", 1024, 768, onReady, onExit)

}

func onReady() {
	systray.SetIcon(getIcon("assets/icon.ico"))
	systray.SetTitle("K8s Switcher")
	systray.SetTooltip("K8s Switcher")

	subMenuClusters := systray.AddMenuItem("Clusters", "K8s Menu")

	var clusterMenu map[string]*systray.MenuItem
	clusterMenu = make(map[string]*systray.MenuItem)

	var namespaceMenu map[string]*systray.MenuItem
	namespaceMenu = make(map[string]*systray.MenuItem)

	for key, element := range clusterMap {

		variable := subMenuClusters.AddSubMenuItem(key, "test")

		for _, namespace := range element {

			namespaceVariable := variable.AddSubMenuItem(namespace, "Descrption")

			namespaceMenu[namespace] = namespaceVariable
		}

		clusterMenu[key] = variable

	}

	systray.AddSeparator()

	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	for key, value := range namespaceMenu {

		go func(key string, value *systray.MenuItem) {
			for {
				select {
				case <-value.ClickedCh:
					if value.Checked() == false {
						value.Check()
						fmt.Println("key: ", key)
						fmt.Println("value: ", value)
						for notCurrentkey, notCurrentValue := range namespaceMenu {
							if notCurrentkey != key {
								notCurrentValue.Uncheck()
							}

						}

					}
					//	return
				}
			}
		}(key, value)
	}
}
