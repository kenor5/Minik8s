package client

import (
	"context"
	"encoding/json"
	"log"
	"minik8s/entity"
	pb "minik8s/pkg/proto"
)

/**
** Kubelet作为客户端给Api Server发请求
**/
func RegisterNode(c pb.ApiServerKubeletServiceClient, in *pb.RegisterNodeRequest) error {
	ctx := context.Background()

	// 调用服务端 RegisterNode 并获取响应
	reply, err := c.RegisterNode(ctx, in)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply.Status)

	return err
}

func UpdatePodStatus(c pb.ApiServerKubeletServiceClient, pod *entity.Pod) error {
	ctx := context.Background()
	podByte, err := json.Marshal(pod)
	log.Println("Begin update Pod Status", string(podByte))
	if err != nil {
		log.Println("parse pod error")
		return err
	}

	updatePodStatusRequest := &pb.UpdatePodStatusRequest{
		Data: podByte,
	}
	// 调用服务端 UpdatePodStatus 并获取响应
	reply, err := c.UpdatePodStatus(ctx, updatePodStatusRequest)
	if err != nil {
		log.Fatal("UpdatePodStatus err!=nil", err)
	}
	log.Println(reply.Status)

	return err
}
