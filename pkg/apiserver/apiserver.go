package apiserver

import (
	"encoding/json"

	// "minik8s/configs"

	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"

	"minik8s/entity"
	"minik8s/pkg/apiserver/ControllerManager/NodeController"
	"minik8s/pkg/apiserver/client"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
)

/**************************************************************************
************************    API Server 主结构    ****************************
***************************************************************************/
type ApiServer struct {
	// conn pb.KubeletApiServerServiceClient
	NodeManager NodeController.NodeController
}

var apiServer *ApiServer

func newApiServer() *ApiServer {
	newServer := &ApiServer{
		NodeManager: *NodeController.NewNodeController(),
	}
	return newServer
}

func ApiServerObject() *ApiServer {
	if apiServer == nil {
		apiServer = newApiServer()
	}
	return apiServer
}

func (master *ApiServer) ApplyPod(in *pb.ApplyPodRequest) (*pb.StatusResponse, error) {
	// 调度(获取conn)
	conn := master.NodeManager.RoundRobin()
	// 发送消息给Kubelet
	err := client.KubeletCreatePod(conn, in)
	if err != nil {
		log.PrintE(err)
		return &pb.StatusResponse{Status: -1}, err
	}

	return &pb.StatusResponse{Status: 0}, err
}

func (master *ApiServer) DeletePod(in *pb.DeletePodRequest) (*pb.StatusResponse, error) {
	//查询Pod对应的Node信息并获取conn
	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	if (in.Data == nil || pod.Status.Phase == entity.Succeed) {
        return &pb.StatusResponse{Status: 0}, err
	}
	// 根据Pod所在的节点的NodeName获得对应的grpc Conn
    conn := master.NodeManager.GetNodeConnByName(pod.Spec.NodeName)
	if conn == nil {
		panic("UnKnown NodeName!\n")
	}
	// 发送消息给Kubelet
	err = client.KubeletDeletePod(conn, in)
	if err != nil {
		log.PrintE(err)
	
		return &pb.StatusResponse{Status: -1}, err
	}

	return &pb.StatusResponse{Status: 0}, err
}

func (master *ApiServer) CreateService(in *pb.ApplyServiceRequest2) (*pb.StatusResponse, error) {
	LivingNodes := master.NodeManager.GetAllLivingNodes()

	for _, node := range LivingNodes {
		// 发送消息给Kubelet
		conn := master.NodeManager.GetNodeConnByName(node.Name)
	    err := client.KubeLetCreateService(conn, in)
	    if err != nil {
			log.PrintE(err)
		    return &pb.StatusResponse{Status: -1}, err
	    }
	}
	return &pb.StatusResponse{Status: 0}, nil
}