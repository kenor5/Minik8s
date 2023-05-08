package main

import (
	"fmt"
	"minik8s/entity"
	"minik8s/tools/yamlParser"
	"time"

	docker "minik8s/pkg/kubelet/container/containerfunc"
)

/**
* Example of creating and destroying a single container
* 1. create a new container according to name and image
* 2. start the container
* 3. stop the container
* 4. remove the container
 */
func main() {
	yamlPath := "./test/pod2.yaml"
	pod := &entity.Pod{}
	b, _ := yamlParser.ParseYaml(pod, yamlPath)
	if !b {
		fmt.Println("test ParseYaml error")
	}
	for _, container := range pod.Spec.Containers {
		id := docker.CreateContainer(container)
		fmt.Printf("container %s created\n", id)
		fmt.Printf("It will be started in 1s\n")
		time.Sleep(time.Second * 1)
		docker.StartContainer(id)
		fmt.Printf("after 10s, container %s will be stopped\n", id)
		time.Sleep(time.Second * 10)
		docker.StopContainer(id)
		fmt.Printf("after 3s, container %s will be removed\n", id)
		time.Sleep(time.Second * 3)
		id, err := docker.RemoveContainer(id)
		if err == nil {
			fmt.Printf("container %s is removed\n", id)
		}
	}
	//fmt.Printf("container %s created\n", id)
	//fmt.Printf("It will be started in 1s\n")
	//time.Sleep(time.Second * 1)
	//docker.StartContainer(id)
	//fmt.Printf("after 40s, container %s will be stopped\n", id)
	//time.Sleep(time.Second * 40)
	//docker.StopContainer(id)
	//fmt.Printf("after 3s, container %s will be removed\n", id)
	//time.Sleep(time.Second * 3)

}
