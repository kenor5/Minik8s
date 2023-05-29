package PodManager

import (
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/kubelet/pod/podfunc"
	"minik8s/tools/log"
	"sync"
)

// PodManager reference to : https://github.dev/kubernetes/kubernetes
// 主要用于管理Pod的MetaData
type PodManager interface {
	// GetPods returns the regular pods bound to the kubelet and their spec.
	GetPods() []*entity.Pod
	// GetPodByFullName returns the (non-mirror) pod that matches full name, as well as
	// whether the pod was found.
	GetPodByFullName(podFullName string) (*entity.Pod, bool)

	// GetPodByName provides the (non-mirror) pod that matches namespace and
	// name, as well as whether the pod was found.
	GetPodByName(namespace string, name string) (*entity.Pod, bool)
	// GetPodByUID provides the (non-mirror) pod that matches pod UID, as well as
	// whether the pod is found.
	GetPodByUID(string) (*entity.Pod, bool)

	AddPod(pod *entity.Pod)
	// UpdatePod updates the given pod in the manager.
	UpdatePod(pod *entity.Pod)
	// DeletePod deletes the given pod from the manager.  For mirror pods,
	// this means deleting the mappings related to mirror pods.  For non-
	// mirror pods, this means deleting from indexes for all non-mirror pods.
	DeletePod(pod *entity.Pod)

	// AddContainerToPod DeletePodByName(name string)
	AddContainerToPod(containerId string, pod *entity.Pod)

	GetContainersByPod(pod *entity.Pod) []string
	// DeleteContainersByPod Delete containerIds in Pod
	DeleteContainersByPod(pod *entity.Pod)
}

// basicManager is a functional Manager.
//
// All fields in basicManager are read-only and are updated calling SetPods,
// AddPod, UpdatePod, or DeletePod.
type basicManager struct {
	// Protects all internal maps.
	lock sync.RWMutex

	// Regular pods indexed by UID.
	podByUID map[string]*entity.Pod

	// Pods indexed by full name for easy access.
	podByFullName   map[string]*entity.Pod
	ContainersByPod map[string][]string
}

// NewPodManager returns a functional Manager.
func NewPodManager() PodManager {
	return &basicManager{
		podByFullName:   map[string]*entity.Pod{},
		podByUID:        map[string]*entity.Pod{},
		ContainersByPod: map[string][]string{},
	}
}

func (pm *basicManager) AddPod(pod *entity.Pod) {
	log.PrintS("a")
	pm.lock.Lock()
	defer pm.lock.Unlock()
	if _, ok := pm.GetPodByName(pod.Metadata.Namespace, pod.Metadata.Name); ok {
		fmt.Printf("pod %v already exist\n", pod.Metadata.Name)
	}
	log.PrintS("e")
	fullName := pod.Metadata.Namespace + pod.Metadata.Name
	pm.podByFullName[fullName] = pod
	log.PrintS("f")
	pm.podByUID[pod.Metadata.Uid] = pod
	log.PrintS("g")
}

func (pm *basicManager) UpdatePod(pod *entity.Pod) {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	fullName := pod.Metadata.Namespace + pod.Metadata.Name
	pm.podByFullName[fullName] = pod

	if pod.Metadata.Uid == "" {
		fmt.Println("Uid not set")
	} else {
		pm.podByUID[pod.Metadata.Uid] = pod
	}
}

func (pm *basicManager) DeletePod(pod *entity.Pod) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	podFullName := pod.Metadata.Namespace + pod.Metadata.Name

	delete(pm.podByUID, pod.Metadata.Uid)
	delete(pm.podByFullName, podFullName)

}

func (pm *basicManager) GetPods() []*entity.Pod {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pods := make([]*entity.Pod, 0, len(pm.podByFullName))
	for _, pod := range pm.podByFullName {
		pods = append(pods, pod)
	}
	return pods
}

func (pm *basicManager) GetPodByUID(uid string) (*entity.Pod, bool) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pod, ok := pm.podByUID[uid]
	return pod, ok
}

func (pm *basicManager) GetPodByName(namespace string, name string) (*entity.Pod, bool) {
	podFullName := namespace + name
	return pm.GetPodByFullName(podFullName)
}

func (pm *basicManager) GetPodByFullName(podFullName string) (*entity.Pod, bool) {
	// log.PrintS("b")
	// pm.lock.RLock()
	// log.PrintS("c")
	// defer pm.lock.RUnlock()
	pod, ok := pm.podByFullName[podFullName]
	log.PrintS("d")
	return pod, ok
}

func (pm *basicManager) AddContainerToPod(containerId string, pod *entity.Pod) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	fullname := pod.Metadata.Namespace + pod.Metadata.Name
	pm.ContainersByPod[fullname] = append(pm.ContainersByPod[fullname], containerId)
}

func (pm *basicManager) GetContainersByPod(pod *entity.Pod) []string {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	fullname := pod.Metadata.Namespace + pod.Metadata.Name
	return pm.ContainersByPod[fullname]
}

func (pm *basicManager) DeleteContainersByPod(pod *entity.Pod) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	fullname := pod.Metadata.Namespace + pod.Metadata.Name
	//删除其中的Container
	podfunc.DeletePod(pm.ContainersByPod[fullname])
	pm.ContainersByPod[fullname] = pm.ContainersByPod[fullname][:0]
	return
}
