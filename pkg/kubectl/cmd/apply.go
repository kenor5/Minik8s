package cmd

import (
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/kubectl/utils"
	"minik8s/tools/log"
	"minik8s/tools/yamlParser"
	"strings"

	"github.com/spf13/cobra"
)

var (
	filename string
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "kubectl apply is to create object",
		Long: `kubectl apply is to create object
		   for example:
		   kubectl apply -f *.yaml`,
		Run: doApply,
	}
)

func init() {
	applyCmd.Flags().StringVarP(&filename, "filename", "f", "", "yaml name")
	err := applyCmd.MarkFlagRequired("filename")
	if err != nil {
		log.LOG("required filename")
		return
	}
}

func doApply(cmd *cobra.Command, args []string) {
	var (
		dirname                  string
		filenameWithoutExtention string
	)

	b, err := yamlParser.FileExists(filename)
	if !b || err != nil {
		log.LOG("file does not exist")
		return
	}

	arr := strings.Split(filename, "/")
	fmt.Println(arr)
	for i := 0; i < len(arr)-1; i++ {
		dirname = dirname + arr[i] + "/"
	}
	if len(dirname) == 0 {
		dirname = "."
	}

	filenameWithoutExtention = strings.Split(arr[len(arr)-1], ".")[0]

	obj, err := utils.GetField(dirname, filenameWithoutExtention, "kind")
	if err != nil {
		log.LOG("file has no such field")
	}

	fmt.Println(dirname, filenameWithoutExtention)

	switch obj {
	case "Pod":
	case "pod":
		pod := &entity.Pod{}
		_, err := yamlParser.ParseYaml(pod, filename)
		if err != nil {
			log.LOG("parse pod failed")
			return
		}
		fmt.Println(pod)
		//url := "http://localhost:8080/pod"
		//err = http.ApplyToServer(url, pod)
		//if err != nil {
		//	return
		//}

	case "Deployment":
	case "deployment":
		deploy := &entity.Deployment{}
		_, err := yamlParser.ParseYaml(deploy, filename)
		if err != nil {
			log.LOG("parse deploy failed")
			return
		}

	case "Service":
	case "service":

	case "Node":
	case "node":
	default:
		log.LOG("there is no object named ")

	}
}
