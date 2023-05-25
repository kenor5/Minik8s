package scale

import "minik8s/entity"

var resourceController *ResourceController

type ResourceController struct {
	PodsName map[string]*entity.Pod
	//AverageCPUUsage    map[string]float64
	//AverageMemoryUsage map[string]uint64
	metricsManager *MetricsManager
}

func (rc *ResourceController) ResourceControllerObjetc() *ResourceController {
	if resourceController == nil {
		resourceController = rc.NewresourceController()
		return resourceController
	} else {
		return resourceController
	}
}

func (rc *ResourceController) GetPodCPUUsage(pod *entity.Pod) (float64, error) {
	//TODO implement me
	//panic("implement me")
	CPUUsage, err := rc.metricsManager.PodCPUUsage(pod)
	if err != nil {
		return 0, err
	}
	return CPUUsage, err
}

func (rc *ResourceController) GetPodMemoryUsage(pod *entity.Pod) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (rc *ResourceController) NewresourceController() *ResourceController {
	//if resourceController==nil{
	//	resourceController=
	//}
	return &ResourceController{
		PodsName: map[string]*entity.Pod{},
		//AverageMemoryUsage: map[string]uint64{},
		//AverageCPUUsage:    map[string]float64{},
		metricsManager: NewMetricsManager(),
	}
}
