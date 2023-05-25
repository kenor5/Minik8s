package main

import (
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/apiserver/scale"
	"minik8s/tools/yamlParser"
)

// GetMetricsTest 查询一个Pod的CPU和Memory使用情况
func GetMetricsTest() {
	var metricsManager *scale.MetricsManager
	metricsManager = scale.NewMetricsManager()

	pod := &entity.Pod{}
	_, err := yamlParser.ParseYaml(pod, "/home/zhaoxi/go/src/minik8s/test/pod2.yaml")
	if err != nil {
		return
	}

	usage, err := metricsManager.PodCPUUsage(pod)
	if err != nil {
		return
	}
	fmt.Printf("%s CPU Usage:%f", pod.Metadata.Name, usage)
	memoryUsage, err := metricsManager.PodMemoryUsage(pod)
	if err != nil {
		return
	}
	fmt.Printf("%s Memory Usage:%d bytes", pod.Metadata.Name, memoryUsage)
}
