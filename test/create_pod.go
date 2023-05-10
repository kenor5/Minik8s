package main

import (
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/kubelet/pod/podfunc"
	"minik8s/tools/yamlParser"
	"time"
)

var yamlPath = "test/pod3.yaml"

func main() {
	// parse yaml
	newPod := &entity.Pod{}
	yamlParser.ParseYaml(newPod, yamlPath)
	fmt.Println("****ParseYaml Pod*****")
	fmt.Println(newPod)
	fmt.Println(newPod.Spec.Containers[0].Resources.Limit["cpu"])
	fmt.Println(newPod.Spec.Containers[0].Resources.Limit["Cpu"])
	ContainerIDs, err := podfunc.CreatePod(newPod)
	fmt.Println("ContainerIDs: ", ContainerIDs)
	if err != nil {
		fmt.Printf("something wrong\n")
	}
	time.Sleep(time.Second * 10)
	fmt.Println("delete all Container")
	podfunc.DeletePod(ContainerIDs)

}
