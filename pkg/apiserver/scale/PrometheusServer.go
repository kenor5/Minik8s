package scale

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
)

const prometheusName = "prom/prometheus"
const prometheusPort = 9090

// 部署启动一个prometheus服务
func StartPrometheusServer() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()
	if err != nil {
		panic(err)
	} // Pull Prometheus image
	reader, err := cli.ImagePull(ctx, "prom/prometheus", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	//fmt.Println("Pulling Prometheus image...")
	exposedPort := nat.Port(fmt.Sprintf("%d/tcp", prometheusPort))
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: prometheusName,
		ExposedPorts: nat.PortSet{
			exposedPort: struct{}{},
		},
	}, &container.HostConfig{
		Privileged: true,
		PortBindings: nat.PortMap{
			exposedPort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprint(prometheusPort),
				},
			},
		},
	}, nil, nil, prometheusName)

	err = cli.ContainerStart(context.Background(), resp.ID, dockertypes.ContainerStartOptions{})
	if err != nil {
		log.Printf("fail to start cadvisor container: %v", err)
		return err
	}
	return err

}
