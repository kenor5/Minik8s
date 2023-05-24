package ScaleController

import (
	"encoding/json"
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/apiserver/ControllerManager"
	"minik8s/pkg/apiserver/scale"
	"minik8s/tools/etcdctl"
	"minik8s/tools/log"
	"strconv"
	"time"
)

const (
	scaleIntervalTime = 30
)

type Controller interface {
	// CreateAutoscaler 新增一个扩缩策略
	CreateAutoscaler(autoscaler *entity.HorizontalPodAutoscaler) error
	// DescribeAutoscalers 打印当前扩缩策略
	DescribeAutoscalers(all bool, autoscalerNames []string) ([]*entity.HorizontalPodAutoscaler, []string)
	//TODO 删除deploument时删除相关策略
}

/*
扩缩的策略实现是，根据etcd的HPA/name中对应的策略，manager查询对应的deployment的整体pod资源使用情况，并且更新currentReplicas字段
根据扩缩策略修改desiredReplicas策略，修改etcd中的Pod内容，通过Pod更新机制通知kubelet对应新增或者减少Pod
*/

// AutoscalerManager
type AutoscalerManager struct {
	metricsManager *scale.MetricsManager
	autoscalers    map[string]*entity.HorizontalPodAutoscaler //
}

func (AM *AutoscalerManager) CreateAutoscaler(autoscaler *entity.HorizontalPodAutoscaler) error {
	autoscalerName := autoscaler.Metadata.Name
	_, ok := AM.autoscalers[autoscalerName]
	if ok {
		log.Printf("HPA %s already exists", autoscalerName)
		return nil
	}
	AM.autoscalers[autoscalerName] = autoscaler
	return nil
}

// 按照scaleInterval指定的时间间隔（默认30s）执行策略，更新状态
func (AM *AutoscalerManager) startAutoscalerMonitor(autoscaler *entity.HorizontalPodAutoscaler) {
	deploymentName := autoscaler.Spec.ScaleTargetRef.Name
	monitorInterval := time.Second * time.Duration(autoscaler.Spec.ScaleInterval)
	ticker := time.NewTicker(monitorInterval) //每30s使用
	for range ticker.C {
		deploymentdata, _ := etcdctl.EtcdGet("Deployment/" + deploymentName)
		if len(deploymentdata.Kvs) == 0 {
			// Deployment 已经被删除
			delete(AM.autoscalers, autoscaler.Metadata.Name)
			return
		}
		deployment := &entity.Deployment{}
		err := json.Unmarshal(deploymentdata.Kvs[0].Value, deployment)
		if err != nil {
			log.PrintE("[startAutoscalerMonitor]Unmarshal deployment error")
			return
		}
		AM.monitorAndScaleDeployment(autoscaler, deployment)
		//重新写入deployment
		deploymentData, _ := json.Marshal(deployment)
		err = etcdctl.EtcdPut("Deployment/"+deployment.Metadata.Name, string(deploymentData))
		if err != nil {
			return
		}
	}
}

//func MemoryParse(Mstr string) string {
//
//}

// 依次查询deployment中pod的资源使用情况
// 容器命名：deploymentname(nginx-deployment)+templet对应HASH(9594276)+PodUID后五位+镜像名称
func (AM *AutoscalerManager) monitorAndScaleDeployment(autoscaler *entity.HorizontalPodAutoscaler, deployment *entity.Deployment) {
	//获取当前deployment的所有Pod
	//var pods []entity.Pod
	pods := ControllerManager.GetPodsBydeployment(deployment.Metadata.Name)
	if pods == nil || pods[0].Metadata.Name == "" {
		panic("[monitorAndScaleDeployment]read Pods error")
	}
	NowPodNums := len(pods)
	if NowPodNums == 0 {
		log.PrintW("[AutoscalerManager] GetPodsBydeployment == 0")
		return
	}
	//已经超过最大运行的MaxReplicas，直接返回
	if autoscaler.Status.CurrentReplicas+autoscaler.Status.DesiredReplicas >= autoscaler.Spec.MaxReplicas {
		log.Print("[AutoscalerManager] Have reached MaxReplicas")
		return
	}

	var deploymentCPUUsage float64
	var cpuUsagePerPod float64
	var cpuUsageAvgPod float64
	var deploymentMemoryUsage uint64
	var memoryUsagePerPod uint64
	var memoryUsageAvgPod uint64
	var err error
	//查询Pod的资源使用情况，计算deployment总资源使用情况
	autoscaler.Status.ObservedGeneration++
	for _, metric := range autoscaler.Spec.Metrics {
		switch metric.Resource.Name {
		case "cpu":
			for _, pod := range pods {
				cpuUsagePerPod, err = AM.metricsManager.PodCPUUsage(&pod)
				if err != nil {
					log.PrintW("[monitorAndScaleDeployment]Get cpuUsage fail", pod.Metadata.Name)
				}
				deploymentCPUUsage += cpuUsagePerPod * 100 //百分比*100
			}
			cpuUsageAvgPod = deploymentCPUUsage / float64(NowPodNums)
			autoscaler.Status.CurrentMetrics[0].ResourceStatus.Current.AverageUtilization = fmt.Sprintf("%f", cpuUsageAvgPod)
			//当平均CPU使用率大于averageUtilization,则增加replica数量
			TargetCPUAvg, _ := strconv.ParseFloat(metric.Resource.Target.AverageUtilization, 64)
			//newReplica := autoscaler.Status.CurrentReplicas + autoscaler.Status.DesiredReplicas
			newReplica := deployment.Spec.Replicas
			if cpuUsageAvgPod > TargetCPUAvg && newReplica < autoscaler.Spec.MaxReplicas {
				//autoscaler.Status.DesiredReplicas++
				deployment.Spec.Replicas++
			} else if newReplica >= autoscaler.Spec.MaxReplicas && cpuUsageAvgPod > TargetCPUAvg {
				//已经超过最大运行的MaxReplicas，直接返回

				fmt.Printf("AUTOSCALER [%s]:Have reached MaxReplicas %d",
					autoscaler.Metadata.Name,
					autoscaler.Spec.MaxReplicas,
				)
				return
			}

			//当平均CPU使用率小于averageUtilization,则减少replica数量，不低于MinReplica
			if cpuUsageAvgPod < TargetCPUAvg && newReplica > autoscaler.Spec.MinReplicas {
				//autoscaler.Status.DesiredReplicas--
				deployment.Spec.Replicas--
			} else if newReplica <= autoscaler.Spec.MaxReplicas && cpuUsageAvgPod < TargetCPUAvg {
				//已经低于最小运行的MInReplicas，直接返回
				fmt.Printf("AUTOSCALER [%s]:Have reached MinReplicas %d",
					autoscaler.Metadata.Name,
					autoscaler.Spec.MinReplicas,
				)
				return
			}
			autoscaler.Status.CurrentReplicas = deployment.Spec.Replicas
			autoscaler.Status.LastScaleTime = time.Now()
			fmt.Printf("AUTOSCALER [%s]: cpu usage per pod reaches %TargetCPUAvg, deployment %s scales out to %d replicas\n",
				autoscaler.Metadata.Name,
				cpuUsagePerPod,
				deployment.Metadata.Name,
				deployment.Spec.Replicas)

		case "memory":
			deploymentMemoryUsage = 0
			for _, pod := range pods {
				memoryUsagePerPod, err = AM.metricsManager.PodMemoryUsage(&pod)
				if err != nil {
					log.PrintW("[monitorAndScaleDeployment]Get cpuUsage fail", pod.Metadata.Name)
				}
				deploymentMemoryUsage += memoryUsagePerPod
			}
			memoryUsageAvgPod = deploymentMemoryUsage / uint64(NowPodNums)
			TargetMemoryAvg, _ := strconv.ParseUint(metric.Resource.Target.AverageUtilization, 10, 64)
			autoscaler.Status.CurrentMetrics[0].ResourceStatus.Current.AverageUtilization = fmt.Sprintf("%d", TargetMemoryAvg)

			newReplica := deployment.Spec.Replicas

			if memoryUsageAvgPod > TargetMemoryAvg && newReplica < autoscaler.Spec.MaxReplicas {
				//autoscaler.Status.DesiredReplicas++
				deployment.Spec.Replicas++
			} else if newReplica >= autoscaler.Spec.MaxReplicas && memoryUsageAvgPod > TargetMemoryAvg {
				//已经低于大运行的MaxReplicas，直接返回
				fmt.Printf("AUTOSCALER [%s]:Have reached MaxReplicas %d",
					autoscaler.Metadata.Name,
					autoscaler.Spec.MinReplicas,
				)
				return
			}

			if memoryUsageAvgPod < TargetMemoryAvg && newReplica > autoscaler.Spec.MinReplicas {
				//autoscaler.Status.DesiredReplicas--
				deployment.Spec.Replicas--
			} else if newReplica <= autoscaler.Spec.MinReplicas && memoryUsageAvgPod < TargetMemoryAvg {
				//已经低于最小运行的MInReplicas，直接返回

				fmt.Printf("AUTOSCALER [%s]:Have reached MinReplicas %d",
					autoscaler.Metadata.Name,
					autoscaler.Spec.MinReplicas,
				)
				return
			}
			autoscaler.Status.CurrentReplicas = deployment.Spec.Replicas
			autoscaler.Status.LastScaleTime = time.Now()
			fmt.Printf("AUTOSCALER [%s]: memory usage per pod reaches %TargetMemoryAvg, deployment %s scales out to %d replicas\n",
				autoscaler.Metadata.Name,
				TargetMemoryAvg,
				deployment.Metadata.Name,
				deployment.Spec.Replicas)

		default:
			log.PrintW("[monitorAndScaleDeployment]Don't support resource:", metric.Resource.Name)
		}

	}

}
