package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/**
* remove a container according to the cotainerID
 */
func RemoveContainer(containerID string) (string, error) {
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()

	err := cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{})

	if err != nil {
		panic(err)
	} else {
		return containerID, err
	}
}
