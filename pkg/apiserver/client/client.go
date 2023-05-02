package client

import (
	"context"
	"log"
	pb "minik8s/pkg/proto"
)

/**
** ApiServer作为客户端给Kubelet发请求
**/
func KubeletCreatePod(c pb.KubeletApiServerServiceClient, in *pb.ApplyPodRequest) error {
	ctx := context.Background()

	// 调用服务端 SimpleRPC 并获取响应
	reply, err := c.CreatePod(ctx, in)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply.Status)

	return err
}
