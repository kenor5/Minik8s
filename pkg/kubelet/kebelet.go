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
	PodManger PodManager.Manager
}

var kubelet Kubelet

// newKubelet creates a new Kubelet object.
func newKubelet() *Kubelet {
	return &Kubelet{
		PodManger: PodManager.NewPodManager(),
	}
}

func KubeletObject() *Kubelet {
	kubelet := newKubelet()
	return kubelet
}

func (kl *Kubelet) CreatePod(pod *entity.Pod) error {
	err := podfunc.CreatePod(pod)
	return err
}

func (kl *Kubelet) GetPods() ([]*entity.Pod, error) {
	pm := kl.PodManger.GetPods()
	return pm, nil
}

func (kl *Kubelet) AddPod(pod *entity.Pod) error {
	//更新元数据
	kl.PodManger.AddPod(pod)
	pod.Status.Phase = entity.Running

	//启动沙箱容器和pod.spec.containers中的容器
	if err := podfunc.CreatePod(pod); err != nil {
		pod.Status.Phase = entity.Failed
		return err
	}

	return nil
}

func (kl *Kubelet) GetPodByName(namespace string, name string) (*entity.Pod, bool) {
	pm, ok := kl.PodManger.GetPodByName(namespace, name)
	return pm, ok
}

func (kl *Kubelet) DeletePod(pod *entity.Pod) {
	kl.PodManger.DeletePod(pod)

}
