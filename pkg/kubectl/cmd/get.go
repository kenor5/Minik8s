package cmd

import (
	"github.com/spf13/cobra"
	"minik8s/tools/log"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get resources info",
	Long: `get resources info
		   for example:
		   kubectl get pod [podname] [-n namespace]  	get pod info of namespace
		   kubectl get deployment [deployment name] [-n namespace] 	get deployment info`,
	Args: cobra.MinimumNArgs(1),
	Run:  doGet,
}

func doGet(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.LOG("get err must have 2 args")
		return
	}
	name := args[1]
	switch args[0] {
	case "po":
	case "pod":
	case "pods":
		getPod(name)
	case "node":
	case "nodes":
		getNode(name)
	case "service":
		getService(name)
	case "function":
		getFunction(name)
	case "deployment":
	case "deploy":
		getDeployment(name)
	}
}

func getPod(name string) {

}

func getNode(name string) {

}

func getDeployment(name string) {

}

func getFunction(name string) {

}

func getService(name string) {

}
