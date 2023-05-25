package containerfunc

import (
	"context"
	"minik8s/tools/log"

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
		log.PrintS("container ", containerID, " start successfully")
	}
}
