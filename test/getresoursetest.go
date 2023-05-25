package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Stats struct {
	CpuStats struct {
		CpuUsage struct {
			TotalUsage int64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemUsage int64 `json:"system_cpu_usage"`
	} `json:"cpu_stats"`
	MemoryStats struct {
		Usage int64 `json:"usage"`
	} `json:"memory_stats"`
}

func GetResourseTest() {

	containerId := "0afb74cb19c1bd687b01f306438172dee4b9a3dfc315bc2a57039e2bd166834a"
	url := fmt.Sprintf("http://localhost:8080/api/v1.3/docker/%s", containerId)
	//url := fmt.Sprintf("http://localhost:8080/docker/%s", containerId)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var stats Stats
	json.NewDecoder(resp.Body).Decode(&stats)
	statsData, _ := json.Marshal(stats)
	println(string(statsData))
	cpuPercent := float64(stats.CpuStats.CpuUsage.TotalUsage) / float64(stats.CpuStats.SystemUsage) * 100.0
	memoryUsage := float64(stats.MemoryStats.Usage) / 1024.0 / 1024.0

	fmt.Printf("CPU: %.2f%%\nMemory: %.2f MB\n", cpuPercent, memoryUsage)
}
