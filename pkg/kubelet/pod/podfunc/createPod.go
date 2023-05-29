package podfunc

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-units"
	"regexp"
	"strconv"

	"minik8s/entity"
	docker "minik8s/pkg/kubelet/container/containerfunc"
	"minik8s/tools/log"
	UUID "minik8s/tools/uuid"
	"time"
)

func CreatePod(pod *entity.Pod) ([]string, error) {
	// Create and Start Pause Container
	log.Print("create pause container\n")

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
		log.Printf("create common container: %s\n", con.Name)
		docker.EnsureImage(con.Image)

		//TODO：待明确Pod中将哪一个目录供Container挂载使用 emptydir?
		//增加卷volume绑定
		//映射示例  Binds: []string{"/path/on/host:/path/in/container:rw"},
		vBinds := make([]string, 0, len(con.VolumeMounts))
		if len(pod.Spec.Volumes) == 0 {
			for _, m := range con.VolumeMounts {
				PodPath := pauseName + pod.Metadata.Uid + "_" + con.Name
				vBinds = append(vBinds, fmt.Sprintf("%v:%v", PodPath, m.MountPath))
			}
		} else {
			for _, m := range con.VolumeMounts {
				for _, pm := range pod.Spec.Volumes {
					if pm.Name == m.Name {
						vBinds = append(vBinds, fmt.Sprintf("%v:%v", pm.HostPath, m.MountPath))
					}
				}
			}
		}

		//增加容器CPU资源限制
		MemoryLimit, err := parseMemorySize(con.Resources.Limit.Memory)
		if err != nil {
			log.PrintE(err)
		}
		CPULimit, err := parseCPUSize(con.Resources.Limit.Cpu)
		if err != nil {
			log.PrintE(err)
		}

		// 容器的限制：128MB 的内存，相当于 134217728 字节，和 0.1 个 CPU 核
		resources := container.Resources{
			Memory:   int64(MemoryLimit),
			NanoCPUs: int64(CPULimit * 1000000000), // 1000000000相当于 1 个 CPU 核,分配0.1
		}

		config := &container.Config{
			Image: con.Image,
			Cmd:   con.Command,
		}
		config.Tty = true //保持容器運行
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
			log.PrintE("ContainerCreate error")
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
	fmt.Printf("Kubelet create Pod and begin Update Status\n")
	containerIP := containerJSON.NetworkSettings.IPAddress
	pod.Status.StartTime = time.Now()
	pod.Status.PodIp = containerIP
	pod.Status.Phase = entity.Running

	log.Printf("Create Pod success! Pod IP: %s\n", containerIP)
	return ContainerIDMap, nil
}

func parseMemorySize(s string) (int, error) {
	re := regexp.MustCompile(`^(\d+)([kKmMgGtTpPeE]i?)?$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid size string: %s", s)
	}
	size, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}
	if len(matches) >= 3 {
		suffix := matches[2]
		switch suffix {
		case "ki", "Ki":
			size *= units.KiB
		case "mi", "Mi":
			size *= units.MiB
		case "gi", "Gi":
			size *= units.GiB
		case "ti", "Ti":
			size *= units.TiB
		case "pi", "Pi":
			size *= units.PiB
		case "e", "E":
			size *= 1024 * units.PiB
		}
	}
	return size, nil
}

func parseCPUSize(cpu string) (float64, error) {
	re := regexp.MustCompile(`^(\d+)m$`)
	match := re.FindStringSubmatch(cpu)
	if len(match) != 2 {
		return 0, fmt.Errorf("invalid cpu %s", cpu)
	}
	n, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, err
	}
	return n / 1000, nil

}
