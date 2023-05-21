package scale

import (
	"context"
	"fmt"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
	"minik8s/pkg/kubelet/container/containerfunc"
)

const (
	prometheusConatinerName = "prometheus"
	prometheusName          = "prom/prometheus"
	prometheusPort          = 9090
	prometheusConfig        = "prometheus.yml"
	ConfigPath              = "/home/zhaoxi/go/src/minik8s/configs/"
	PrometheusConfigPath    = "/etc/prometheus/"
)

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
	err = containerfunc.EnsureImage("prom/prometheus")
	if err != nil {
		return err
	}
	vBinds := make([]string, 0)
	vBinds = append(vBinds, fmt.Sprintf("%v:%v", ConfigPath, PrometheusConfigPath))
	//fmt.Println("Pulling Prometheus image...")
	exposedPort := nat.Port(fmt.Sprintf("%d/tcp", prometheusPort))
	resp, err := cli.ContainerCreate(ctx, &container.Config{
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
		Binds: vBinds,
	}, nil, nil, prometheusConatinerName)

	err = cli.ContainerStart(context.Background(), resp.ID, dockertypes.ContainerStartOptions{})
	if err != nil {
		log.Printf("fail to start prometheus container: %v", err)
		return err
	}
	return err

}
