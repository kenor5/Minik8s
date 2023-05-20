package containerfunc

import (
	"context"
	"minik8s/tools/log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

/**
* stop a container according to the cotainerID
 */
func StopContainer(containerID string) {

	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()

	timeout := int(time.Second * 10)
	stopOpts := container.StopOptions{Timeout: &timeout}

	err := cli.ContainerStop(context.Background(), containerID, stopOpts)

	if err != nil {
		panic(err)
	} else {
		log.Print("container %s is stopped\n", containerID)
	}
}
