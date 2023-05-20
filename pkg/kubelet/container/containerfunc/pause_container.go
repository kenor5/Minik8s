package containerfunc

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
	"minik8s/entity"
)

func CreatePauseContainer(pod *entity.Pod) (string, string, error) {
	fmt.Printf("**********start create pause container***********\n")
	// Step1: 保证镜像存在
	EnsureImage(entity.PauseImage)

	// Step2: 暴露Ports,选择未分配的空闲端口分配
	// 因为所有容器与pause container共享相同的网络命名空间
	fmt.Printf("Populate exposed ports\n")
	//ports := make(map[nat.Port]struct{})
	exportsPort, portMap := generatePorts(pod.Spec.Containers)
	/*for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			ports[nat.Port(fmt.Sprintf("%v/tcp", port.ContainerPort))] = struct{}{}
		}
	}*/

	// Step3: 创建pause container
	fmt.Printf("create container\n")
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()

	pauseName := pod.Metadata.Name + "-" + "pauseContainer"

	body, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        entity.PauseImage,
		ExposedPorts: exportsPort,
	}, &container.HostConfig{
		IpcMode:      "shareable",
		PortBindings: portMap,
	}, nil, nil, pauseName)
	if err != nil {
		println(err)
		fmt.Printf("start pause_container failed!")
	} else {
		fmt.Printf("start pause_container success!")
	}
	StartContainer(body.ID)

	return body.ID, pauseName, err
}

func generatePorts(cons []entity.Container) (nat.PortSet, nat.PortMap) {
	exportPorts := make(nat.PortSet)
	portMap := make(nat.PortMap)
	for _, con := range cons {
		portsMappings := con.Ports
		for _, ports := range portsMappings {
			if ports.ContainerPort != "" {
				port, err := nat.NewPort(ports.Protocol, ports.ContainerPort)
				if err != nil {
					log.Fatal(err)
				}
				exportPorts[port] = struct{}{}

				if ports.HostPort != "" {
					portBind := nat.PortBinding{HostPort: ports.HostPort}
					tmp := make([]nat.PortBinding, 0, 1)
					tmp = append(tmp, portBind)
					portMap[port] = tmp
				}
			}
		}
	}

	return exportPorts, portMap
}
