package podfunc

import (
	"context"
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/kubelet/docker"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func CreatePod(pod *entity.Pod) error {
	// Create and Start Pause Container
	fmt.Printf("create pause container\n")
	pauseContainerId, err := docker.CreatePauseContainer(pod)
	if err != nil {
		return err
	}

	// Create and Start Common Container
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()

	containerMode := "container:" + pauseContainerId

	for _, con := range pod.Spec.Containers {
		fmt.Printf("create common container: %s\n", con.Name)
		docker.EnsureImage(con.Image)

		config := &container.Config{
			Image: con.Image,
		}

		HostConfig := &container.HostConfig{
			PidMode:     container.PidMode(containerMode),
			IpcMode:     container.IpcMode(containerMode),
			NetworkMode: container.NetworkMode(containerMode),
		}

		containerName := pod.Metadata.Name + "_" + con.Name

		body, err := cli.ContainerCreate(context.Background(), config, HostConfig, nil, nil, containerName)
		if err != nil {
			return err
		}

		docker.StartContainer(body.ID)
	}

	// 更新Pod.Status
	// Inspect the container
	containerJSON, err := cli.ContainerInspect(context.Background(), pauseContainerId)
	if err != nil {
		panic(err)
	}

	// Get the container's IP address
	containerIP := containerJSON.NetworkSettings.IPAddress
	pod.Status.PodIp = containerIP
	pod.Status.Phase = "Running"
	// TODO:给Kubelet分配真正的IP
	pod.Status.HostIp = "127.0.0.1"

	return nil
}
