package cmd

import (
	"github.com/spf13/cobra"
	"minik8s/tools/log"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete resources cluster",
	Long: `delete resources cluster
		   for example:
		   kubectl delete deployment [deployment name] [-n namespace] 	delete deployment`,
	Run: doDelete,
}

func doDelete(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.LOG("describe err must have 2 args")
		return
	}
	name := args[1]
	switch args[0] {
	case "po":
	case "pod":
	case "pods":
		deletePod(name)
	case "node":
	case "nodes":
		deleteNode(name)
	case "service":
		deleteService(name)
	case "function":
		deleteFunction(name)
	case "deployment":
	case "deploy":
		deleteDeployment(name)
	}
}

func deletePod(name string) {

}

func deleteNode(name string) {

}

func deleteDeployment(name string) {

}

func deleteFunction(name string) {

}

func deleteService(name string) {

}
