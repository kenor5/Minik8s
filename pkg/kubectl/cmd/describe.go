package cmd

import (
	"github.com/spf13/cobra"
	"minik8s/tools/log"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "show detailed info of object",
	Long: `show detailed info of object
		   for example:
		   kubectl describe pod [pod name] [-n namespace]  get detailed info of pod`,
	Run: doDescribe,
}

func doDescribe(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.PrintE("describe err must have 2 args")
		return
	}
	name := args[1]
	switch args[0] {
	case "po","pod","pods":
		describePod(name)
	case "node","nodes":
		describeNode(name)
	case "service":
		describeService(name)
	case "function":
		describeFunction(name)
	case "deployment","deploy":
		describeDeployment(name)
	}
}

func describePod(name string) {

}

func describeNode(name string) {

}

func describeDeployment(name string) {

}

func describeFunction(name string) {

}

func describeService(name string) {

}
