package client

import (
	"context"
	"minik8s/tools/log"
	pb "minik8s/pkg/proto"
)

/**
** ApiServer作为客户端给Kubelet发请求
**/
func KubeletCreatePod(c pb.KubeletApiServerServiceClient, in *pb.ApplyPodRequest) error {
	ctx := context.Background()

	// 调用服务端 SimpleRPC 并获取响应
	_, err := c.CreatePod(ctx, in)
	if err != nil {
		log.PrintE(err)
	}

	return err
}

func KubeletDeletePod(c pb.KubeletApiServerServiceClient, in *pb.DeletePodRequest) error {
	ctx := context.Background()

	// 调用服务端 SimpleRPC 并获取响应
	_, err := c.DeletePod(ctx, in)
	if err != nil {
		log.PrintE(err)
	}

	return err
}


func KubeLetCreateService(c pb.KubeletApiServerServiceClient, in *pb.ApplyServiceRequest2) error {
	ctx := context.Background()
	if c == nil {
		log.PrintE("client is nil")
	}
	_, err := c.CreateService(ctx, in)
	if err != nil {
		log.PrintE(err)
	}

	return err
}	

func KubeLetDeleteService(c pb.KubeletApiServerServiceClient, in *pb.DeleteServiceRequest2) error {
	ctx := context.Background()

	_, err := c.DeleteService(ctx, in)
	if err != nil {
		log.PrintE(err)
	}

	return err
}	

func KubeLetCreateDns(c pb.KubeletApiServerServiceClient, in *pb.ApplyDnsRequest) error {
	ctx := context.Background()
	if c == nil {
		log.PrintE("client is nil")
	}
	_, err := c.CreateDns(ctx, in)
	if err != nil {
		log.PrintE(err)
	}
	
	return err
}	

func KubeLetDeleteDns(c pb.KubeletApiServerServiceClient, in *pb.DeleteDnsRequest) error {
	ctx := context.Background()

	_, err := c.DeleteDns(ctx, in)
	if err != nil {
		log.PrintE(err)
	}

	return err
}	