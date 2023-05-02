package main

import (
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/kubelet/pod/podfunc"
	"minik8s/tools/yamlParser"
)

var yamlPath = "./pod2.yaml"

func main() {
	// parse yaml
	newPod := &entity.Pod{}
	yamlParser.ParseYaml(newPod, yamlPath)
	fmt.Println(newPod)

	err := podfunc.CreatePod(newPod)
	if err != nil {
		fmt.Printf("something wrong")
	}
}
