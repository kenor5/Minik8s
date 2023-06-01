package client

import (
	"context"
	"encoding/json"
	"minik8s/configs"
	"minik8s/entity"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
	"minik8s/tools/yamlParser"
)

/**
** Kubelet作为客户端给Api Server发请求  in *pb.RegisterNodeRequest
**/
func RegisterNode(c pb.ApiServerKubeletServiceClient, hostName string, hostIp string) ([][]byte, error) {
	// 获取主机信息
	newNode := &entity.Node{}
	yamlParser.ParseYaml(newNode, "/home/zhaoxi/go/src/minik8s/configs/node/node1.yaml")
	newNode.Ip = hostIp
	newNode.KubeletUrl = hostIp + configs.KubeletGrpcPort
	nodeByte, err := json.Marshal(newNode)
	if err != nil {
		log.PrintE("parse node error")
		return nil, err
	}

	ctx := context.Background()
	// 组装消息
	in := &pb.RegisterNodeRequest{
		NodeData: nodeByte,
	}
	// 调用服务端 RegisterNode 并获取响应
	reply, err := c.RegisterNode(ctx, in)
	if err != nil {
		log.PrintE(err)
	}
	log.Print(reply.PodData)

	return reply.PodData, err
}

func UpdatePodStatus(c pb.ApiServerKubeletServiceClient, pod *entity.Pod) error {
	ctx := context.Background()
	podByte, err := json.Marshal(pod)
	log.Print("Begin update Pod Status", string(podByte))
	if err != nil {
		log.Print("parse pod error")
		return err
	}

	updatePodStatusRequest := &pb.UpdatePodStatusRequest{
		Data: podByte,
	}
	// 调用服务端 UpdatePodStatus 并获取响应
	reply, err := c.UpdatePodStatus(ctx, updatePodStatusRequest)
	if err != nil {
		log.PrintE("UpdatePodStatus err!=nil", err)
	}
	log.Print(reply.Status)

	return err
}
