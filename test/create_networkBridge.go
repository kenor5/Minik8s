package main

import (
	"context"
	"fmt"

	"minik8s/pkg/kubelet/netbridge"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main_1() {
	FlannelNetInterfaceName := "flannel.1"
	// 获取IP地址
	flannelIP, _ := netbridge.GetNetInterfaceIPv4Addr(FlannelNetInterfaceName)

	fmt.Println("flannelIP: ", flannelIP)

	// 判断flannel_bridge网桥是否存在，如果已存在，则不需要再创建；否则创建网桥
	exist, networkID := netbridge.FindNetworkBridge()
	if !exist {
		networkID, _ = netbridge.CreateNetBridge(flannelIP)
	}

	fmt.Printf("networkID: %s\n", networkID)

	// 测试用，打印出所有的网桥信息
	ShowNetWorkBridge()
}

// 打印出所有的网桥信息
func ShowNetWorkBridge() {
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	defer cli.Close()
	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	for _, network := range networks {
		fmt.Printf("network name:%s, network ID: %s, network IP: %s\n", network.Name, network.ID, network.IPAM.Config)
	}
}
