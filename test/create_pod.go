package main

import (
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/kubelet/pod/podfunc"
	"minik8s/tools/yamlParser"
)

var yamlPath = "test/pod3.yaml"

func main() {
	// parse yaml
	newPod := &entity.Pod{}
	yamlParser.ParseYaml(newPod, yamlPath)
	fmt.Println(newPod)

	ContainerIDs, err := podfunc.CreatePod(newPod)
	fmt.Println("ContainerIDs: ", ContainerIDs)
	if err != nil {
		fmt.Printf("something wrong")
	}
}
