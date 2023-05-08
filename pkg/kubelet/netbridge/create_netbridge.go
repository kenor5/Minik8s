package netbridge

import (
	"context"
	"fmt"
	"net"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func CreateNetBridge(flannelIP string) (string, error) {
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	defer cli.Close()

	_, subnet, _ := net.ParseCIDR(flannelIP + "/24")
	gatewayIP := make(net.IP, len(subnet.IP))
	copy(gatewayIP, subnet.IP)
	gatewayIP[3] = 1

	fmt.Printf("subnet:%s, gatewayIP:%s", subnet.String(), gatewayIP.String())
	ipamConfig := network.IPAMConfig{
		Subnet:  subnet.String(),
		Gateway: gatewayIP.String(),
	}

	networkCreateOpetion := types.NetworkCreate{
		Driver: "bridge",
		IPAM: &network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{ipamConfig},
		},
		Options: map[string]string{
			"com.docker.network.bridge.name": "flannel_bridge",
		},
	}

	newNetwork, err := cli.NetworkCreate(
		context.Background(),
		"flannel_bridge",
		networkCreateOpetion,
	)

	fmt.Printf("Netbridge %s has been created\n", newNetwork.ID)

	return newNetwork.ID, err
}

func FindNetworkBridge() (bool, string) {
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
		if network.Name == "flannel_bridge" {
			fmt.Println("flannel_bridge already exist, no need to create")
			return true, network.ID
		}
	}

	return false, ""
}
