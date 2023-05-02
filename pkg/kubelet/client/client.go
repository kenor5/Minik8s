package client

import (
	"context"
	"log"
	pb "minik8s/pkg/proto"
)

/**
** Kubelet作为客户端给Api Server发请求
**/
func RegisterNode(c pb.ApiServerKubeletServiceClient, in *pb.RegisterNodeRequest) error {
	ctx := context.Background()

	// 调用服务端 SimpleRPC 并获取响应
	reply, err := c.RegisterNode(ctx, in)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply.Status)

	return err
}
