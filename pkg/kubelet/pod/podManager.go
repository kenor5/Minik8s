package pod

import (
	"fmt"
	"minik8s/entity"
	"sync"

)

// reference to : https://github.dev/kubernetes/kubernetes


type Manager interface {
	// GetPods returns the regular pods bound to the kubelet and their spec.
	GetPods() []*entity.Pod
	// GetPodByFullName returns the (non-mirror) pod that matches full name, as well as
	// whether the pod was found.
	GetPodByFullName(podFullName string) (*entity.Pod, bool)
	// GetPodByName provides the (non-mirror) pod that matches namespace and
	// name, as well as whether the pod was found.
	GetPodByName(namespace, name string) (*entity.Pod, bool)
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
	podByFullName       map[string]*entity.Pod

}


// NewBasicPodManager returns a functional Manager.
func NewBasicPodManager() Manager {
	pm := &basicManager{}
	pm.lock.Lock()
	defer pm.lock.Unlock()

	pm.podByUID = map[string]*entity.Pod{}
	pm.podByFullName = map[string]*entity.Pod{}
	
	return pm
}


func (pm *basicManager) AddPod(pod *entity.Pod) {
	pm.UpdatePod(pod)
}

func (pm *basicManager) UpdatePod(pod *entity.Pod) {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	fullName := pod.Metadata.Namespace + pod.Metadata.Name
	pm.podByFullName[fullName] = pod

	if pod.Metadata.Uid == ""{
		fmt.Println("Uid not set")
	}else {
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
	pods := make([]*entity.Pod, 0, len(pm.podByUID))
	for _,po := range pm.podByUID {
		pods = append(pods, po)
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
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pod, ok := pm.podByFullName[podFullName]
	return pod, ok
}
