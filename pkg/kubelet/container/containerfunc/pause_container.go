package containerfunc

import (
	"context"
	"fmt"
	"minik8s/entity"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func CreatePauseContainer(pod *entity.Pod, NetBridgeID string) (string, error) {
	fmt.Printf("**********start create pause container***********\n")
	// Step1: 保证镜像存在
	EnsureImage(entity.PauseImage)

	// Step2: 暴露Ports
	// 因为所有容器与pause container共享相同的网络命名空间
	fmt.Printf("Populate exposed ports\n")
	ports := make(map[nat.Port]struct{})
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			ports[nat.Port(fmt.Sprintf("%v/tcp", port.ContainerPort))] = struct{}{}
		}
	}

	// Step3: 创建pause container
	fmt.Printf("create container\n")
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()

	// 加入网络配置，使用Flannel网桥，保证IP分配
	endpointConfig := &network.EndpointSettings{
		NetworkID: NetBridgeID,
	}
	// "flannel_bridge"为网桥名称
	networkconfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"flannel_bridge": endpointConfig,
		},
	}

	pauseName := pod.Metadata.Name + "_" + "pauseContainer"

	body, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        entity.PauseImage,
		ExposedPorts: ports,
	}, &container.HostConfig{
		IpcMode: "shareable",
	}, networkconfig, nil, pauseName)

	fmt.Printf("start container %s\n", err)
	StartContainer(body.ID)

	return body.ID, err
}
