package resourse

import (
	"context"
	"fmt"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"minik8s/tools/log"
)

// const cAdvisorImage = "google/cadvisor:v0.36.0"
const cAdvisorImage = "google/cadvisor:latest"
const cadvisorPort = 8080
const cAdvisorName = "kubeletcAdvisor"

func StartcAdvisor() error {
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()
	_, err := cli.ImagePull(context.Background(), cAdvisorImage, dockertypes.ImagePullOptions{})
	if err != nil {
		log.PrintE("Can't pull image:", cAdvisorImage)
		panic(err)
	}

	//根据官网要求设置映射https://github.com/google/cadvisor
	vBinds := []string{
		"/:/rootfs:ro",
		"/var/run:/var/run:ro",
		"/sys:/sys:ro",
		"/var/lib/docker/:/var/lib/docker:ro",
		"/dev/disk/:/dev/disk:ro",
	}
	exposedPort := nat.Port(fmt.Sprintf("%d/tcp", cadvisorPort))

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: cAdvisorImage,
		ExposedPorts: nat.PortSet{
			exposedPort: struct{}{},
		},
		Cmd: strslice.StrSlice{"--max_housekeeping_interval=2s"},
	}, &container.HostConfig{
		Binds:      vBinds,
		Privileged: true,
		PortBindings: nat.PortMap{
			exposedPort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprint(cadvisorPort),
				},
			},
		},
	}, nil, nil, cAdvisorName)
	if err != nil {
		return err
	}
	err = cli.ContainerStart(context.Background(), resp.ID, dockertypes.ContainerStartOptions{})
	if err != nil {
		log.Print("fail to start cadvisor container", err)
		return err
	}
	log.Print("starts cadvisor container success!")
	return err
}
