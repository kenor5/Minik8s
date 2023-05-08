package main

import (
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/kubelet/netbridge"
	"minik8s/pkg/kubelet/pod/podfunc"
	"minik8s/tools/yamlParser"
)

var yamlPath = "test/pod3.yaml"

func main_1() {
	// parse yaml
	newPod := &entity.Pod{}
	yamlParser.ParseYaml(newPod, yamlPath)
	fmt.Println(newPod)

	FlannelNetInterfaceName := "flannel.1"
	// 获取IP地址
	flannelIP, _ := netbridge.GetNetInterfaceIPv4Addr(FlannelNetInterfaceName)

	fmt.Println("flannelIP: ", flannelIP)

	// 判断flannel_bridge网桥是否存在，如果已存在，则不需要再创建；否则创建网桥
	exist, netBridgeID := netbridge.FindNetworkBridge()
	if !exist {
		netBridgeID, _ = netbridge.CreateNetBridge(flannelIP)
	}

	ContainerIDs, err := podfunc.CreatePod(newPod, netBridgeID)
	fmt.Println("ContainerIDs: ", ContainerIDs)
	if err != nil {
		fmt.Printf("something wrong")
	}
}
