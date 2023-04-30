package apiserver

import (
	"context"
	"fmt"
	"log"
	"minik8s/configs"
	"minik8s/tools/etcdctl"
	"net"

	pb "minik8s/proto"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type server struct {
	// 继承 protoc-gen-go-grpc 生成的服务端代码
	pb.UnimplementedApiServerKubeletServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Println("[Api Server] Kubelet call sayHello...")

	log.Println(in)
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

/*
	referenct to:
	https://blog.csdn.net/qq_43580193/article/details/127577709
*/

func Run() {
	// 开启etcd
	cli, err := etcdctl.Start(configs.EtcdStartPath)
	if err != nil {
		return
	}
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			fmt.Println("etcd close error")
		}
	}(cli)

	// 注册请求处理接口
	listen, err := net.Listen("tcp", configs.GrpcPort)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建gRPC服务器
	svr := grpc.NewServer()
	// 将实现的接口注册进 gRPC 服务器
	pb.RegisterApiServerKubeletServiceServer(svr, &server{})
	log.Println("Apiserver For Kubelet gRPC Server starts running...")

	// 启动 gRPC 服务器
	err = svr.Serve(listen)
	if err != nil {
		log.Fatal(err)
		return
	}
}
