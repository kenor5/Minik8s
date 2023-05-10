package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"minik8s/configs"
	"minik8s/entity"
	"minik8s/pkg/kubelet"
	pb "minik8s/pkg/proto"
	"net"

	"google.golang.org/grpc"
)

/**
** Kubelet作为服务端接受来自Api Server的请求
**/
type server struct {
	// 继承 protoc-gen-go-grpc 生成的服务端代码
	pb.UnimplementedKubeletApiServerServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Println("[Kubelet] Api Server call sayHello...")

	log.Println(in)
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

func (s *server) CreatePod(ctx context.Context, in *pb.ApplyPodRequest) (*pb.StatusResponse, error) {
	log.Println("[Kubelet] Api Server call Create Pod...")
	log.Println(in)

	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	if err != nil {
		fmt.Println("pod unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}

	// 调用kubelet对象实现创建Pod的真正逻辑
	err = kubelet.KubeletObject().CreatePod(pod)

	if err != nil {
		fmt.Println("create pod err")
		return &pb.StatusResponse{Status: -1}, err
	}
	fmt.Println("[Kubelet] Create Pod Success")
	return &pb.StatusResponse{Status: 0}, err
}

func (s *server) DeletePod(ctx context.Context, in *pb.DeletePodRequest) (*pb.StatusResponse, error) {
	log.Println("[Kubelet] Api Server call delete Pod...")
	log.Println(in)

	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	if err != nil {
		fmt.Println("pod unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}

	// 调用kubelet对象实现删除 Pod的真正逻辑
	err = kubelet.KubeletObject().DeletePod(pod)

	if err != nil {
		fmt.Println("delete pod err")
		return &pb.StatusResponse{Status: -1}, err
	}
	fmt.Println("[Kubelet] delete Pod Success")
	return &pb.StatusResponse{Status: 0}, err
}

/*********************************************************
**********************  Kubelet主程序   *************************
**********************************************************/
func Run() {
	println("[kubelet] running...")

	/**
	 *    Kubelet启动时向APIServer注册
	 **/
	kubelet.KubeletObject().RegisterNode()
	println("[kubelet] has registered to apiserver...")

	/**
	 *    Kubelet启动自己的服务端，接受来自ApiServer的消息
	**/
	listen, err := net.Listen("tcp", configs.KubeletGrpcPort)
	if err != nil {
		fmt.Println(err)
	}

	// 创建gRPC服务器
	svr := grpc.NewServer()
	// 将实现的接口注册进 gRPC 服务器
	pb.RegisterKubeletApiServerServiceServer(svr, &server{})

	// 启动 gRPC 服务器
	err = svr.Serve(listen)
	if err != nil {
		log.Fatal(err)
		return
	}
}
