package main

import (
	"fmt"
	"minik8s/pkg/apiserver/ControllerManager"
)

func main_3() {
	Labels := map[string]string{
		"app": "myApp",
	}
	selectedPods := ControllerManager.GetPodsByLabels(&Labels)

	// 遍历列表并查看元素
	for element := selectedPods.Front(); element != nil; element = element.Next() {
		value := element.Value
		fmt.Println(value)
	}
}
