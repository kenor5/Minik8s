package kubelet

import (
	// "context"
	// "encoding/json"
	// "fmt"
	// "log"
	// "minik8s/configs"
	"minik8s/entity"
	"minik8s/pkg/kubelet/pod/podfunc"

	// pb "minik8s/pkg/proto"
	// "net"

	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
	"minik8s/pkg/kubelet/pod/PodManager"
)

/**************************************************************************
************************    Kubelet主结构    *******************************
***************************************************************************/
type Kubelet struct {
	PodManger *PodManager.Manager
}

var kubelet *Kubelet

// newKubelet creates a new Kubelet object.
func newKubelet() *Kubelet {
	return &Kubelet{}
}

func KubeletObject() *Kubelet {
	if kubelet == nil {
		kubelet = newKubelet()
	}
	return kubelet
}

func (kl *Kubelet) CreatePod(pod *entity.Pod) error {
	podfunc.CreatePod(pod)
	return nil
}

func (kl *Kubelet) GetPods() ([]*entity.Pod, error) {
	return nil, nil
}
