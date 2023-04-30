package kubelet

import (
	"context"
	"log"
	"minik8s/configs"
	pb "minik8s/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sayHello(c pb.ApiServerKubeletServiceClient) {
	ctx := context.Background()

	// 调用服务端 SimpleRPC 并获取响应
	reply, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Kubelet"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply.GetReply())
}

func Run() {
	println("[kubelet] running...")

	// 连接服务端，因为我们没有SSL证书，因此这里需要禁用安全传输
	apiserver_url := "127.0.0.1" + configs.GrpcPort
	dial, err := grpc.Dial(apiserver_url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer dial.Close()

	conn := pb.NewApiServerKubeletServiceClient(dial)
	sayHello(conn)
}
