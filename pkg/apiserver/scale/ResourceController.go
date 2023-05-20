package scale

import "minik8s/entity"

var resourceController *ResourceController

type MetricsManager interface {
	// PodCPUUsage queries the average CPU usage of a given pod in the past certain seconds.
	PodCPUUsage(pod *entity.Pod) (float64, error)
	// PodMemoryUsage queries the average memory usage of a given pod in the past certain seconds.
	PodMemoryUsage(pod *entity.Pod) (uint64, error)
}

type ResourceController struct {
	PodsName           map[string]*entity.Pod
	AverageCPUUsage    map[string]float64
	AverageMemoryUsage map[string]uint64
}

func (rc *ResourceController) ResourceControllerObjetc() *ResourceController {
	if resourceController == nil {
		resourceController = rc.NewresourceController()
		return resourceController
	} else {
		return resourceController
	}
}

func (rc *ResourceController) PodCPUUsage(pod *entity.Pod) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (rc *ResourceController) PodMemoryUsage(pod *entity.Pod) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (rc *ResourceController) NewresourceController() *ResourceController {
	//if resourceController==nil{
	//	resourceController=
	//}
	return &ResourceController{
		PodsName:           map[string]*entity.Pod{},
		AverageMemoryUsage: map[string]uint64{},
		AverageCPUUsage:    map[string]float64{},
	}
}
