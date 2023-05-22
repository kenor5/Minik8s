package ScaleController

import (
	"encoding/json"
	"minik8s/entity"
	"minik8s/pkg/apiserver/scale"
	"minik8s/tools/etcdctl"
	"minik8s/tools/log"
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
		log.Prinf("HPA %s already exists", autoscalerName)
		return nil
	}
	AM.autoscalers[autoscalerName] = autoscaler
	return nil
}

// 按照scaleInterval指定的时间间隔（默认30s）执行策略，更新状态
func (AM *AutoscalerManager) startAutoscalerMonitor(autoscaler *entity.HorizontalPodAutoscaler) {
	deploymentName := autoscaler.Spec.ScaleTargetRef.Name
	monitorInterval := time.Second * time.Duration(autoscaler.Spec.ScaleInterval)
	ticker := time.NewTicker(monitorInterval)
	for range ticker.C {
		deploymentdata, _ := etcdctl.EtcdGet(deploymentName)
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
	}
}

// 依次查询deployment中pod的资源使用情况
// 容器命名：deploymentname(nginx-deployment)+templet对应HASH(9594276)+PodUID后五位+镜像名称
func (AM *AutoscalerManager) monitorAndScaleDeployment(autoscaler *entity.HorizontalPodAutoscaler, deployment *entity.Deployment) {
	//获取当前deployment的所有Pod
	//var pods []*entity.Pod

}
