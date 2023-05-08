package kubelet

import (
	// "context"
	// "encoding/json"
	// "fmt"

	"log"
	"minik8s/configs"
	"minik8s/entity"
	"minik8s/pkg/kubelet/pod/podfunc"

	pb "minik8s/pkg/proto"
	// "net"

	"minik8s/pkg/kubelet/client"
	"minik8s/pkg/kubelet/container/ContainerManager"
	"minik8s/pkg/kubelet/pod/PodManager"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/**************************************************************************
************************    Kubelet主结构    *******************************
***************************************************************************/
type Kubelet struct {
	connToApiServer  pb.ApiServerKubeletServiceClient // kubelet连接到apiserver的conn
	podManger        *PodManager.Manager
	containerManager *ContainerManager.ContainerManager
}

var kubelet *Kubelet

// newKubelet creates a new Kubelet object.
func newKubelet() *Kubelet {
	newKubelet := &Kubelet{
		containerManager: ContainerManager.NewContainerManager(),
	}
	apiserver_url := "127.0.0.1" + configs.GrpcPort
	newKubelet.connToApiServer, _ = ConnectToApiServer(apiserver_url)
	return newKubelet
}

func KubeletObject() *Kubelet {
	if kubelet == nil {
		kubelet = newKubelet()
	}
	return kubelet
}

func (kl *Kubelet) CreatePod(pod *entity.Pod) error {
	// 实际创建Pod,IP等信息在这里更新进Pod.Status中
	ContainerIds, err := podfunc.CreatePod(pod)
	if err != nil {
		return err
	}

	// 维护ContainerRuntimeManager
	kubelet.containerManager.SetContainerIDsByPodName(pod, ContainerIds)

	// 更新PodStatus
	client.UpdatePodStatus(kubelet.connToApiServer, pod)
	return nil
}

func (kl *Kubelet) DeletePod(pod *entity.Pod) error {
	// 获取Pod中所有的ContainerId并且删除该映射
	containerIds := kubelet.containerManager.GetContainerIDsByPodName(pod.Metadata.Name)
	kubelet.containerManager.DeletePodNameToContainerIds(pod.Metadata.Name)

	// 实际停止并删除Pod中的所有容器
	podfunc.DeletePod(containerIds)

	// 更新Pod的状态
	pod.Status.Phase = entity.Succeed
	client.UpdatePodStatus(kubelet.connToApiServer, pod)
	return nil
}

func (kl *Kubelet) GetPods() ([]*entity.Pod, error) {
	return nil, nil
}

func (kl *Kubelet) RegisterNode() error {
	registerNodeRequest := &pb.RegisterNodeRequest{
		NodeName:   "node1",
		KubeletUrl: "127.0.0.1" + configs.KubeletGrpcPort,
	}
	client.RegisterNode(kubelet.connToApiServer, registerNodeRequest)
	return nil
}

func ConnectToApiServer(apiserver_url string) (pb.ApiServerKubeletServiceClient, error) {
	dial, err := grpc.Dial(apiserver_url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// defer dial.Close()

	conn := pb.NewApiServerKubeletServiceClient(dial)
	return conn, err
}
