package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/getlantern/systray"
)

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}

var clusterMap map[string][]string

func main() {

	clusterMap = make(map[string][]string)

	clusterMap["test-cluster-1"] = []string{"clus1-namespace1", "clus1-namespace2", "clus1-namespace3", "clus1-namespace4"}
	clusterMap["test-cluster-2"] = []string{"clus2-namespace1", "clus2-namespace2", "clus2-namespace3", "clus2-namespace4"}

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
