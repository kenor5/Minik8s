package containerfunc

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/**
* start a container according to the cotainerID
 */
func StartContainer(containerID string) {

	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	defer cli.Close()

	err := cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})

	if err == nil {
		fmt.Println("container", containerID, "start successfully")
	}
}
