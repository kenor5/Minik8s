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
	"minik8s/pkg/kubelet/pod/PodManager"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/**************************************************************************
************************    Kubelet主结构    *******************************
***************************************************************************/
type Kubelet struct {
	connToApiServer pb.ApiServerKubeletServiceClient
	PodManger       *PodManager.Manager
}

var kubelet *Kubelet

// newKubelet creates a new Kubelet object.
func newKubelet() *Kubelet {
	newKubelet := &Kubelet{}
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
	podfunc.CreatePod(pod)

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
