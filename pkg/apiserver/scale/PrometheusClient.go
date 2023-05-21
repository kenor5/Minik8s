package scale

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v3"
	"minik8s/entity"
	"minik8s/pkg/kubelet/container/containerfunc"
	"minik8s/tools/log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	PrometheusAddress    string        = "http://localhost:9090"
	QueryTimeout         time.Duration = 5 * time.Second
	UsageComputeDuration string        = "30s"
)

// MetricsManager monitors the CPU and memory usage of all the ready pods at set intervals.
type MetricsManager interface {
	// PodCPUUsage queries the average CPU usage of a given pod in the past certain seconds.
	PodCPUUsage(pod *entity.Pod) (float64, error)
	// PodMemoryUsage queries the average memory usage of a given pod in the past certain seconds.
	PodMemoryUsage(pod *entity.Pod) (uint64, error)
}

// 初始化API查询CPU和Memory
type metricsManagerInner struct {
	prometheusAPI v1.API
}

func NewMetricsManager() MetricsManager {
	client, err := api.NewClient(api.Config{
		Address: PrometheusAddress,
	})
	if err != nil {
		log.PrintE(err)
	}
	return &metricsManagerInner{
		prometheusAPI: v1.NewAPI(client),
	}
}

func (mm *metricsManagerInner) PodCPUUsage(pod *entity.Pod) (float64, error) {
	var queryBuilder strings.Builder

	// Pause container
	pauseName := pod.Metadata.Name + "_pause"
	containerQuery := containerCPUUsageQuery(pauseName)
	queryBuilder.WriteString(containerQuery)

	// Other containers
	for _, container := range pod.Spec.Containers {
		containerName := pod.Metadata.Name + "_" + container.Name
		containerQuery = containerCPUUsageQuery(containerName)
		queryBuilder.WriteString(" or ")
		queryBuilder.WriteString(containerQuery)
	}

	// 查询总的Pod CPU使用率
	query := "sum(" + queryBuilder.String() + ")"

	// Query Promethus
	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()
	result, warnings, err := mm.prometheusAPI.Query(ctx, query, time.Now())
	if err != nil {
		log.PrintE("fail to get cpu usage from prometheus:", err)
		return 0.0, err
	}
	if len(warnings) > 0 {
		log.PrintE("warnings from prometheus", warnings)
	}
	if result.(model.Vector).Len() == 0 {
		return 0.0, fmt.Errorf("fail to get cpu usage for pod %s: no data from prometheus", pod.Metadata.Name)
	}

	fmt.Printf("pod %s cpu usage: %f\n", pod.Metadata.Name, float64(result.(model.Vector)[0].Value))
	return float64(result.(model.Vector)[0].Value), nil
}

func (mm *metricsManagerInner) PodMemoryUsage(pod *entity.Pod) (uint64, error) {
	var queryBuilder strings.Builder

	// Pause container
	pauseName := pod.Metadata.Name + "_pause"
	containerQuery := containerMemoryUsageQuery(pauseName)
	queryBuilder.WriteString(containerQuery)

	// Other containers
	for _, container := range pod.Spec.Containers {
		containerName := pod.Metadata.Name + "_" + container.Name
		containerQuery = containerMemoryUsageQuery(containerName)
		queryBuilder.WriteString(" or ")
		queryBuilder.WriteString(containerQuery)
	}

	// 查询总的Pod Memory使用率
	query := "sum(" + queryBuilder.String() + ")"

	// Query Promethus
	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()
	result, warnings, err := mm.prometheusAPI.Query(ctx, query, time.Now())
	if err != nil {
		log.PrintE("fail to get memory usage from prometheus:", err)
		return 0, err
	}
	if len(warnings) > 0 {
		log.PrintW("warnings from prometheus:", warnings)
	}
	if result.(model.Vector).Len() == 0 {
		return 0, fmt.Errorf("fail to get memory usage for pod %s", pod.Metadata.Name)
	}

	fmt.Printf("pod %s memory usage: %d bytes\n", pod.Metadata.Name, uint64(result.(model.Vector)[0].Value))
	return uint64(result.(model.Vector)[0].Value), nil
}

// containerCPUUsageQuery 生成 PromQL 查询语句，查询一个容器的CPU过去30s平均使用率
func containerCPUUsageQuery(containerName string) string {
	var query strings.Builder
	query.WriteString("sum(rate(container_cpu_usage_seconds_total{name=\"")
	query.WriteString(containerName)
	query.WriteString("\"}[")
	query.WriteString(UsageComputeDuration)
	query.WriteString("])) by (name)")
	return query.String()
}

// containerMemoryUsageQuery 生成 PromQL 查询语句，查询一个容器的Memory过去30s平均使用率
func containerMemoryUsageQuery(containerName string) string {
	var query strings.Builder
	query.WriteString("avg_over_time(container_memory_usage_bytes{name=\"")
	query.WriteString(containerName)
	query.WriteString("\"}[")
	query.WriteString(UsageComputeDuration)
	query.WriteString("])")
	return query.String()
}

type static_configs struct {
}

type Config struct {
	Global struct {
		ScrapeInterval     string `yaml:"scrape_interval"`
		EvaluationInterval string `yaml:"evaluation_interval"`
	} `yaml:"global"`
	Alerting struct {
		AlertManagers []struct {
			StaticConfigs []struct {
				Targets []string `yaml:"targets"`
			} `yaml:"static_configs"`
		} `yaml:"alertmanagers"`
	} `yaml:"alerting"`
	RuleFiles     []string `yaml:"rule_files"`
	ScrapeConfigs struct {
		job []struct {
			jobName       string `yaml:"job_name"`
			StaticConfigs []struct {
				Targets []string `yaml:"targets"`
			} `yaml:"static_configs"`
		}
	} `yaml:"scrape_configs"`
}

// GeneratePrometheusTargets 使用HostIP和port(9090)注册job到Prometheus配置文件中
func GeneratePrometheusTargets(nodes []*entity.Node) error {
	// 打开配置文件，读取内容
	configFile := ConfigPath + prometheusConfig
	content, err := os.ReadFile(configFile)
	config := &Config{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		fmt.Printf("Failed to read config file: %s", err)
		return err
	}

	// 将 hostIP 和 port 追加到文件尾部
	for _, node := range nodes {
		newTarget := fmt.Sprintf("%s:%s", node.Ip, strconv.Itoa(8080))
		config.ScrapeConfigs.job[1].StaticConfigs[0].Targets = append(config.ScrapeConfigs.job[1].StaticConfigs[0].Targets, newTarget)
	}

	content, _ = yaml.Marshal(config)
	// 将修改后的内容写入文件
	err = os.WriteFile(configFile, content, 0644)
	if err != nil {
		fmt.Printf("Failed to write config file: %s", err)
		return err
	}
	//重启服务
	containerfunc.ReStartContainer(prometheusConatinerName)
	return nil
}
