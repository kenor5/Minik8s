package containerfunc

import (
	"context"
	"fmt"
	"minik8s/entity"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func CreatePauseContainer(pod *entity.Pod) (string, string, error) {
	fmt.Printf("**********start create pause container***********\n")
	// Step1 : Ensure Image
	EnsureImage(entity.PauseImage)

	// Step2: Populate exposed ports
	// Because all the containers share the same network namespace with pause container
	fmt.Printf("Populate exposed ports\n")
	ports := make(map[nat.Port]struct{})
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			ports[nat.Port(fmt.Sprintf("%v/tcp", port.ContainerPort))] = struct{}{}
		}
	}

	// Step3 Create Container
	fmt.Printf("create container\n")
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()

	pauseName := pod.Metadata.Name + "_" + "pauseContainer"

	body, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        entity.PauseImage,
		ExposedPorts: ports,
	}, &container.HostConfig{
		IpcMode: "shareable",
	}, nil, nil, pauseName)

	fmt.Printf("start container %s\n", err)
	StartContainer(body.ID)

	return body.ID, pauseName, err
}
