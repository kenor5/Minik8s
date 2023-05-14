package podfunc

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"minik8s/entity"
	docker "minik8s/pkg/kubelet/container/containerfunc"
	UUID "minik8s/tools/uuid"
)

func CreatePod(pod *entity.Pod) ([]string, error) {
	// Create and Start Pause Container
	fmt.Printf("create pause container\n")

	// 该map返回Pod中的ContainerID
	ContainerIDMap := []string{}

	pauseContainerId, pauseName, err := docker.CreatePauseContainer(pod)
	if err != nil {
		return nil, err
	}

	ContainerIDMap = append(ContainerIDMap, pauseContainerId)

	// Create and Start Common Container
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()

	pauseContainerMode := "container:" + pauseContainerId
	//设定UUID
	pod.Metadata.Uid = UUID.UUID()

	for _, con := range pod.Spec.Containers {
		fmt.Printf("create common container: %s\n", con.Name)
		docker.EnsureImage(con.Image)

		//TODO：待明确Pod中将哪一个目录供Container挂载使用 emptydir?
		//增加卷volume绑定
		//映射示例  Binds: []string{"/path/on/host:/path/in/container:rw"},
		vBinds := make([]string, 0, len(con.VolumeMounts))
		for _, m := range con.VolumeMounts {
			PodPath := pauseName + pod.Metadata.Uid + "_" + con.Name
			vBinds = append(vBinds, fmt.Sprintf("%v:%v", PodPath, m.MountPath))
		}
		//增加容器CPU资源限制
		//resources := container.Resources{}
		fmt.Println(con.Resources.Limit["cpu"])
		//if bytes, ok := con.Resources.Limit["CPU"]; ok {
		//	//if int64(bytes) < 0 {
		//	//	return fmt.Errorf("memory limit overflow: %v", bytes)
		//	//} else {
		//	//	resources, = int64(bytes)
		//	//}
		//}
		// 容器的限制：128MB 的内存，相当于 134217728 字节，和 1 个 CPU 核
		resources := container.Resources{
			Memory:   134217728,
			NanoCPUs: 1000000000, // 相当于 1 个 CPU 核
		}

		config := &container.Config{
			Image: con.Image,
			Cmd:   con.Command,
		}

		HostConfig := &container.HostConfig{
			PidMode:     container.PidMode(pauseContainerMode),
			IpcMode:     container.IpcMode(pauseContainerMode),
			NetworkMode: container.NetworkMode(pauseContainerMode),
			Binds:       vBinds,
			Resources:   resources,
		}

		containerName := pod.Metadata.Name + "_" + con.Name

		body, err := cli.ContainerCreate(context.Background(), config, HostConfig, nil, nil, containerName)
		if err != nil {
			DeletePod(ContainerIDMap)
			return nil, err
		}

		docker.StartContainer(body.ID)

		ContainerIDMap = append(ContainerIDMap, body.ID)
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
	pod.Status.Phase = entity.Running
	// TODO:给Kubelet分配真正的IP
	pod.Status.HostIp = "127.0.0.1"

	fmt.Printf("Create Pod success! Pod IP: %s\n", containerIP)

	return ContainerIDMap, nil
}
