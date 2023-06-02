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
	MetricsManager *scale.MetricsManager
	Autoscalers    map[string]*entity.HorizontalPodAutoscaler //
}

func (AM *AutoscalerManager) CreateAutoscaler(autoscaler *entity.HorizontalPodAutoscaler) error {
	autoscalerName := autoscaler.Metadata.Name
	_, ok := AM.Autoscalers[autoscalerName]
	if ok {
		log.Printf("HPA %s already exists", autoscalerName)
		return nil
	}
	AM.Autoscalers[autoscalerName] = autoscaler
	return nil
}

func (AM *AutoscalerManager) DeleteAutoscaler(autoscalerName string) error {
	_, ok := AM.Autoscalers[autoscalerName]
	if ok {
		log.Printf("HPA %s in Autoscalers", autoscalerName)
		delete(AM.Autoscalers, autoscalerName)
		return nil
	} else {
		log.PrintfW("HPA %s not in Autoscalers", autoscalerName)
		return nil
	}

}

// 按照scaleInterval指定的时间间隔（默认30s）执行策略，更新状态
func (AM *AutoscalerManager) StartAutoscalerMonitor(autoscaler *entity.HorizontalPodAutoscaler) {
	log.PrintS("[AutoscalerMonitor]Begin Monitor " + autoscaler.Metadata.Name)
	deploymentName := autoscaler.Spec.ScaleTargetRef.Name
	monitorInterval := time.Second * time.Duration(autoscaler.Spec.ScaleInterval)
	ticker := time.NewTicker(monitorInterval) //每30s使用
	for range ticker.C {
		// for debug
		log.PrintE("[AutoscalerMonitor]Begin Monitor " + autoscaler.Metadata.Name)

		deploymentdata, _ := etcdctl.EtcdGet("Deployment/" + deploymentName)
		if len(deploymentdata.Kvs) == 0 {
			// Deployment 已经被删除
			delete(AM.Autoscalers, autoscaler.Metadata.Name)
			return
		}
		autoscalerName := autoscaler.Metadata.Name
		_, ok := AM.Autoscalers[autoscalerName]
		if !ok {
			log.PrintfW("[AutoscalerMonitor]%s has been del", autoscalerName)
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
		log.Printf("[AutoscalerMonitor]Put %s :%s", deployment.Metadata.Name, string(deploymentData))
		err = etcdctl.EtcdPut("Deployment/"+deployment.Metadata.Name, string(deploymentData))
		if err != nil {
			log.PrintE("Etcd Put Deployment error")
			return
		}
		HPAData, _ := json.Marshal(autoscaler)
		etcdctl.EtcdPut("HPA/"+autoscaler.Metadata.Name, string(HPAData))
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

	// for debug
	log.PrintE("begin monitorAndScaleDeployment")

	pods := ControllerManager.GetPodsBydeployment(deployment.Metadata.Name)
	if pods == nil || pods[0].Metadata.Name == "" {
		panic("[monitorAndScaleDeployment]read Pods error")
	}
	//NowPodNums := len(pods)
	//if NowPodNums == 0 {
	//	log.PrintW("[AutoscalerManager] GetPodsBydeployment == 0")
	//	return
	//}
	autoscaler.Status.CurrentReplicas = deployment.Status.Replicas
	autoscaler.Status.DesiredReplicas = deployment.Spec.Replicas - deployment.Status.Replicas
	//已经超过最大运行的MaxReplicas，直接返回
	if autoscaler.Status.CurrentReplicas+autoscaler.Status.DesiredReplicas >= autoscaler.Spec.MaxReplicas {
		// for debug
		log.PrintE("Have reached MaxReplicas")

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
	newReplica := deployment.Spec.Replicas

	// for debug
	log.PrintE("在遍历autoscaler.Spec.Metrics循环之前")

	for _, metric := range autoscaler.Spec.Metrics {
		// for debug
		log.PrintE(metric.Resource.Name)

		switch metric.Resource.Name {
		case "cpu":
			for _, pod := range pods {
				cpuUsagePerPod, err = AM.MetricsManager.PodCPUUsage(&pod)
				if err != nil {
					log.PrintW("[monitorAndScaleDeployment]Get cpuUsage fail ", pod.Metadata.Name)
				}
				deploymentCPUUsage += cpuUsagePerPod * 100 //百分比*100

				// for debug
				log.PrintE("cpuUsagePerPod: ", cpuUsageAvgPod, " deploymentCPUUsage: ", deploymentCPUUsage)
			}
			cpuUsageAvgPod = deploymentCPUUsage / float64(newReplica)
			autoscaler.Status.CurrentMetrics[0].ResourceStatus.Current.AverageUtilization = fmt.Sprintf("%f", cpuUsageAvgPod)
			//当平均CPU使用率大于averageUtilization,则增加replica数量


			TargetCPUAvg, _ := strconv.ParseFloat(metric.Resource.Target.AverageUtilization, 64)

			// for debug
			log.PrintE("TargetCPUAvg: ", TargetCPUAvg, " cpuUsageAvgPod: ", cpuUsageAvgPod)

			//newReplica := autoscaler.Status.CurrentReplicas + autoscaler.Status.DesiredReplicas
			//newReplica := deployment.Spec.Replicas
			if cpuUsageAvgPod > TargetCPUAvg*1.05 && newReplica < autoscaler.Spec.MaxReplicas {
				//autoscaler.Status.DesiredReplicas++

				muti := int32(cpuUsageAvgPod / TargetCPUAvg)
				if muti > 1 {
					deployment.Spec.Replicas = muti * newReplica
				} else {
					deployment.Spec.Replicas = newReplica + 1
				}
				log.PrintS(deployment.Metadata.Name + " add replica to " + string(deployment.Spec.Replicas))
			} else if newReplica >= autoscaler.Spec.MaxReplicas && cpuUsageAvgPod > TargetCPUAvg {

			// for debug
			log.PrintE("reach max replica")

				//已经超过最大运行的MaxReplicas，直接返回
				deployment.Spec.Replicas = autoscaler.Spec.MaxReplicas
				log.Printf("AUTOSCALER [%s]:Have reached MaxReplicas %d",
					autoscaler.Metadata.Name,
					autoscaler.Spec.MaxReplicas,
				)
				return
			}

			//当平均CPU使用率小于averageUtilization,则减少replica数量，不低于MinReplica
			if cpuUsageAvgPod < TargetCPUAvg*0.95 && newReplica > autoscaler.Spec.MinReplicas {

				// for debug
				log.PrintE("decrease replica, before:", deployment.Spec.Replicas)
				
				//autoscaler.Status.DesiredReplicas--
				muti := int32(TargetCPUAvg / cpuUsageAvgPod)
				if muti > 1 {
					log.Print(deployment.Metadata.Name + " sub replica " + strconv.Itoa(int(muti)))
					deployment.Spec.Replicas = deployment.Spec.Replicas / muti
				} else {
					log.Print(deployment.Metadata.Name + " sub replica ")
					deployment.Spec.Replicas = newReplica - 1
				}
				log.PrintE("decrease replica, after:", deployment.Spec.Replicas)
			} else if newReplica <= autoscaler.Spec.MaxReplicas && cpuUsageAvgPod < TargetCPUAvg {
				//已经低于最小运行的MInReplicas，直接返回
				log.Printf("AUTOSCALER [%s]:Have reached MinReplicas %d",
					autoscaler.Metadata.Name,
					autoscaler.Spec.MinReplicas,
				)
				deployment.Spec.Replicas = autoscaler.Spec.MinReplicas
			}
			if deployment.Spec.Replicas < autoscaler.Spec.MinReplicas {
				deployment.Spec.Replicas = autoscaler.Spec.MinReplicas
			} else if deployment.Spec.Replicas > autoscaler.Spec.MaxReplicas {
				deployment.Spec.Replicas = autoscaler.Spec.MaxReplicas
			}
			autoscaler.Status.CurrentReplicas = deployment.Spec.Replicas
			autoscaler.Status.DesiredReplicas = deployment.Spec.Replicas - deployment.Status.Replicas
			autoscaler.Status.LastScaleTime = time.Now()
			log.Printf("AUTOSCALER [%s]: cpu usage per pod reaches %f, deployment %s scales out to %d replicas\n",
				autoscaler.Metadata.Name,
				cpuUsageAvgPod,
				deployment.Metadata.Name,
				deployment.Spec.Replicas)

		case "memory":
			deploymentMemoryUsage = 0
			for _, pod := range pods {
				memoryUsagePerPod, err = AM.MetricsManager.PodMemoryUsage(&pod)
				if err != nil {
					log.PrintW("[monitorAndScaleDeployment]Get cpuUsage fail", pod.Metadata.Name)
				}
				deploymentMemoryUsage += memoryUsagePerPod
			}
			memoryUsageAvgPod = deploymentMemoryUsage / uint64(newReplica)
			TargetMemoryAvg, _ := strconv.ParseUint(metric.Resource.Target.AverageUtilization, 10, 64)
			autoscaler.Status.CurrentMetrics[0].ResourceStatus.Current.AverageUtilization = fmt.Sprintf("%d", TargetMemoryAvg)

			newReplica := deployment.Spec.Replicas

			if memoryUsageAvgPod > TargetMemoryAvg && newReplica < autoscaler.Spec.MaxReplicas {
				//autoscaler.Status.DesiredReplicas++
				muti := int32(memoryUsageAvgPod / TargetMemoryAvg)
				if muti > 1 {
					deployment.Spec.Replicas = muti * newReplica
				} else {
					deployment.Spec.Replicas = newReplica + 1
				}
				log.Print(deployment.Metadata.Name + " add replica to " + string(deployment.Spec.Replicas))
				//deployment.Spec.Replicas++
			} else if newReplica >= autoscaler.Spec.MaxReplicas && memoryUsageAvgPod > TargetMemoryAvg {
				//已经低于大运行的MaxReplicas，直接返回
				deployment.Spec.Replicas = autoscaler.Spec.MaxReplicas
				fmt.Printf("AUTOSCALER [%s]:Have reached MaxReplicas %d\n",
					autoscaler.Metadata.Name,
					autoscaler.Spec.MinReplicas,
				)
				return
			}

			if memoryUsageAvgPod < TargetMemoryAvg && newReplica > autoscaler.Spec.MinReplicas {
				//autoscaler.Status.DesiredReplicas--
				//deployment.Spec.Replicas--
				muti := int32(TargetMemoryAvg / memoryUsageAvgPod)
				if muti > 1 {
					log.Print(deployment.Metadata.Name + " sub replica " + strconv.Itoa(int(muti)))
					deployment.Spec.Replicas = deployment.Spec.Replicas / muti
				} else {
					log.Print(deployment.Metadata.Name + " sub 1 replica ")
					deployment.Spec.Replicas = newReplica - 1
				}
			} else if newReplica <= autoscaler.Spec.MinReplicas && memoryUsageAvgPod < TargetMemoryAvg {
				//已经低于最小运行的MInReplicas，直接返回
				deployment.Spec.Replicas = autoscaler.Spec.MinReplicas
				fmt.Printf("AUTOSCALER [%s]:Have reached MinReplicas %d\n",
					autoscaler.Metadata.Name,
					autoscaler.Spec.MinReplicas,
				)
				return
			}
			if deployment.Spec.Replicas < autoscaler.Spec.MinReplicas {
				deployment.Spec.Replicas = autoscaler.Spec.MinReplicas
			} else if deployment.Spec.Replicas > autoscaler.Spec.MaxReplicas {
				deployment.Spec.Replicas = autoscaler.Spec.MaxReplicas
			}
			autoscaler.Status.CurrentReplicas = deployment.Spec.Replicas
			autoscaler.Status.DesiredReplicas = deployment.Spec.Replicas - deployment.Status.Replicas
			autoscaler.Status.LastScaleTime = time.Now()
			fmt.Printf("AUTOSCALER [%s]: memory usage per pod reaches %d, deployment %s scales out to %d replicas\n",
				autoscaler.Metadata.Name,
				TargetMemoryAvg,
				deployment.Metadata.Name,
				deployment.Spec.Replicas)

		default:
			log.PrintW("[monitorAndScaleDeployment]Don't support resource:", metric.Resource.Name)
		}

	}

}
