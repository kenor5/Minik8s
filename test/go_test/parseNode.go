package main

import (
	"fmt"
	"minik8s/entity"
	"minik8s/tools/yamlParser"
)

var NodeyamlPath = "/root/go/src/minik8s/configs/node/node1.yaml"

func main_2() {
	// parse yaml
	newNode := &entity.Node{}
	yamlParser.ParseYaml(newNode, NodeyamlPath)
	fmt.Println("****ParseYaml Pod*****")
	fmt.Println(newNode)
	// fmt.Println(newPod.Spec.Containers[0].Resources.Limit["cpu"])
	// fmt.Println(newPod.Spec.Containers[0].Resources.Limit["Cpu"])
	// ContainerIDs, err := podfunc.CreatePod(newPod)
	// fmt.Println("ContainerIDs: ", ContainerIDs)
	// if err != nil {
	// 	fmt.Printf("something wrong\n")
	// }
	// time.Sleep(time.Second * 10)
	// fmt.Println("delete all Container")
	// podfunc.DeletePod(ContainerIDs)

}