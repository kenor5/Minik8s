package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"minik8s/configs"
	"minik8s/entity"
	"minik8s/pkg/kubelet"
	"minik8s/pkg/kubelet/netbridge"
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
	 *    创建网桥,保证Pod全局唯一的IP分配，和Pod间通信
	 **/
	FlannelNetInterfaceName := "flannel.1"
	// 获取IP地址
	flannelIP, _ := netbridge.GetNetInterfaceIPv4Addr(FlannelNetInterfaceName)
	// 判断flannel_bridge网桥是否存在，如果已存在，则不需要再创建；否则创建网桥
	exist, NetBridgeId := netbridge.FindNetworkBridge()
	if !exist {
		NetBridgeId, _ = netbridge.CreateNetBridge(flannelIP)
	}
	// 存入全局变量kubelet中方便后续使用
	kubelet.KubeletObject().SetMember(flannelIP, NetBridgeId)

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
