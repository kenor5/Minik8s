package scale

import (
	"context"
	"fmt"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"minik8s/configs"
	"minik8s/pkg/kubelet/container/containerfunc"
	"minik8s/tools/log"
)

const (
	prometheusConatinerName = "minik8s_prometheus"
	prometheusName          = "prom/prometheus:latest"
	prometheusPort          = 9090
	prometheusConfig        = "prometheus.yml"

	PrometheusConfigPath = "/etc/prometheus/"
)

// 部署启动一个prometheus服务
func StartPrometheusServer() error {
	// 确定容器是否已经存在，如果存在则不再启动
	exist, err := containerfunc.CheckContainerRunning(prometheusConatinerName)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()
	if err != nil {
		panic(err)
	} // Pull Prometheus image
	err = containerfunc.EnsureImage(prometheusName)
	if err != nil {
		return err
	}
	vBinds := make([]string, 0)
	vBinds = append(vBinds, fmt.Sprintf("%v:%v", configs.PromtheusConfigPath, PrometheusConfigPath))
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
	if err != nil {
		log.PrintfE("fail to create prometheus container", err)
		return err
	}
	err = cli.ContainerStart(context.Background(), resp.ID, dockertypes.ContainerStartOptions{})
	if err != nil {
		log.PrintfE("fail to start prometheus container: %v", err)
		return err
	}
	return err

}
