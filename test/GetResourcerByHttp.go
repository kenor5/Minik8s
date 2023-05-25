package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"log"
	"time"
)

const (
	PrometheusAddress    string        = "http://192.168.245.136:9090"
	QueryTimeout         time.Duration = 5 * time.Second
	UsageComputeDuration string        = "10s"
)

// 获取指定容器的CPU和内存用量
func GetContainerStats(containerName string) (float64, float64, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return 0.0, 0.0, err
	}

	// 获取容器指标
	ctx := context.Background()
	stats, err := cli.ContainerStats(ctx, containerName, false)
	if err != nil {
		return 0.0, 0.0, err
	}
	defer stats.Body.Close()

	// 从HTTP响应中读取指标
	var stream types.StatsJSON
	dec := json.NewDecoder(stats.Body)
	if err := dec.Decode(&stream); err != nil {
		return 0.0, 0.0, err
	}

	// 获取CPU用量（单位:纳秒）
	cpuUsage := stream.CPUStats.CPUUsage.TotalUsage
	preCPUUsage := stream.PreCPUStats.CPUUsage.TotalUsage
	cpuDelta := cpuUsage - preCPUUsage
	sysDelta := stream.CPUStats.SystemUsage - stream.PreCPUStats.SystemUsage
	cpuPercent := float64(cpuDelta) / float64(sysDelta) * float64(len(stream.CPUStats.CPUUsage.PercpuUsage)) * 100.0

	// 获取内存用量（单位：字节）
	memStat := stream.MemoryStats
	memUsage := float64(memStat.Usage - memStat.Stats["cache"])
	return cpuPercent, memUsage, nil
}

// 从Prometheus中获取容器CPU用量
func GetContainerCPUMetrics(containerName string, duration time.Duration) (float64, error) {
	cli, err := api.NewClient(api.Config{
		Address: PrometheusAddress,
	})
	if err != nil {
		return 0.0, err
	}

	v1api := v1.NewAPI(cli)
	query := fmt.Sprintf("avg(rate(container_cpu_usage_seconds_total{name=\"%s\"}[%s])) by (name)", containerName, duration.String())
	//avg(rate(container_cpu_usage_seconds_total{name="fervent_lichterman"}[30s])) by (name)
	result, _, err := v1api.Query(context.Background(), query, time.Now())
	if err != nil {
		return 0.0, err
	}

	if val, ok := result.(model.Vector); ok {
		if len(val) > 0 {
			return float64(val[0].Value), nil
		}
	}

	return 0.0, nil
}

// 从Prometheus中获取容器Memory用量
func GetContainerMemoryMetrics(containerName string, duration time.Duration) (float64, error) {
	cli, err := api.NewClient(api.Config{
		Address: PrometheusAddress,
	})
	if err != nil {
		return 0.0, err
	}

	v1api := v1.NewAPI(cli)
	query := fmt.Sprintf("avg(container_memory_usage_bytes{name=\"%s\"}) by (name)", containerName)
	result, _, err := v1api.Query(context.Background(), query, time.Now())
	if err != nil {
		return 0.0, err
	}

	if val, ok := result.(model.Vector); ok {
		if len(val) > 0 {
			return float64(val[0].Value), nil
		}
	}

	return 0.0, nil
}

func GetResourseByHTTP() {
	containerName := "fervent_lichterman"

	cpuUsage, memUsage, err := GetContainerStats(containerName)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Container CPU Usage: %v\n", cpuUsage)
	fmt.Printf("Container Memory Usage: %v\n", memUsage)

	cpuPercent, err := GetContainerCPUMetrics(containerName, 30*1000000000)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Container CPU Percent: %v\n", cpuPercent)

	memMetrics, err := GetContainerMemoryMetrics(containerName, 30*1000000000)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Container Memory Metrics: %v\n", memMetrics)
}
