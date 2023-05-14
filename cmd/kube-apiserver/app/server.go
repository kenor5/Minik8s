package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"minik8s/configs"

	"minik8s/entity"
	"minik8s/tools/etcdctl"
	"net"

	pb "minik8s/pkg/proto"

	"minik8s/pkg/apiserver"

	clientv3 "go.etcd.io/etcd/client/v3"

	"minik8s/pkg/apiserver/ControllerManager"

	"google.golang.org/grpc"
)

/*
	referenct to:
	https://blog.csdn.net/qq_43580193/article/details/127577709
*/

type server struct {
	// 继承 protoc-gen-go-grpc 生成的服务端代码
	pb.UnimplementedApiServerKubeletServiceServer
	pb.UnimplementedApiServerKubectlServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Println("[Api Server] Kubelet call sayHello...")

	log.Println(in)
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

// ApplyPod 客户端为Kubectl
func (s *server) ApplyPod(ctx context.Context, in *pb.ApplyPodRequest) (*pb.StatusResponse, error) {

	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	if err != nil {
		fmt.Println("pod unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}

	cli, err := etcdctl.NewClient()
	defer cli.Close()
	if err != nil {
		fmt.Println("etcd client connetc error")
	}
	fmt.Println("put etcd", in.Data)
	etcdctl.Put(cli, "Pod/"+pod.Metadata.Name, string(in.Data))

	return apiserver.ApiServerObject().ApplyPod(in)
}

func (s *server) DeletePod(ctx context.Context, in *pb.DeletePodRequest) (*pb.StatusResponse, error) {

	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("connect to etcd error")
	}
	out, err := etcdctl.Get(cli, "Pod/"+string(in.Data))

	if len(out.Kvs) == 0 {
		return apiserver.ApiServerObject().DeletePod(&pb.DeletePodRequest{
			Data: nil,
		})
	} else {
		return apiserver.ApiServerObject().DeletePod(&pb.DeletePodRequest{
			Data: out.Kvs[0].Value,
		})
	}
}

// GetPod TODO: get pods后不跟PodName返回所有的Pod
func (s *server) GetPod(ctx context.Context, in *pb.GetPodRequest) (*pb.GetPodResponse, error) {
	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("connect to etcd error")
	}
	out, err := etcdctl.Get(cli, "Pod/"+string(in.PodName))
	if len(out.Kvs) == 0 {
		return &pb.GetPodResponse{PodData: nil}, nil
	} else {
		return &pb.GetPodResponse{PodData: out.Kvs[0].Value}, nil
	}
}

// RegisterNode 客户端为Kubelet
func (s *server) RegisterNode(ctx context.Context, in *pb.RegisterNodeRequest) (*pb.StatusResponse, error) {

	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("etcd client connetc error")
	}
	fmt.Println("[ApiServer] Regiseter Node: put kubelet_url in etcd", in.KubeletUrl)
	etcdctl.Put(cli, "Node/"+in.NodeName, string(in.KubeletUrl))

	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) UpdatePodStatus(ctx context.Context, in *pb.UpdatePodStatusRequest) (*pb.StatusResponse, error) {
	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	if err != nil {
		fmt.Println("pod unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}

	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("etcd client connetc error")
	}
	fmt.Println("put etcd", in.Data)
	etcdctl.Put(cli, "Pod/"+pod.Metadata.Name, string(in.Data))
	return &pb.StatusResponse{Status: 0}, err
}

// GetService Service
func (s *server) GetService(ctx context.Context, in *pb.GetServiceRequest) (*pb.GetServiceResponse, error) {
	//TODO
	return &pb.GetServiceResponse{Data: nil}, nil
}

func (s *server) DeleteService(ctx context.Context, in *pb.DeleteServiceRequest) (*pb.StatusResponse, error) {
	//TODO
	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) ApplyService(ctx context.Context, in *pb.ApplyServiceRequest) (*pb.StatusResponse, error) {
	service := &entity.Service{}
	err := json.Unmarshal(in.Data, service)
	if err != nil {
		return &pb.StatusResponse{Status: -1}, err
	}

	// 放进etcd
	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("etcd client connetc error")
	}
	fmt.Println("put etcd", in.Data)
	etcdctl.Put(cli, "Service/"+service.Metadata.Name, string(in.Data))

	// 获取符合条件的Pod
	selectedPods := ControllerManager.GetPodsByLabels(&service.Metadata.Labels)
	ControllerManager.PrintList(selectedPods)

	return &pb.StatusResponse{Status: 0}, nil
}

// Deployment
func (s *server) GetDeployment(ctx context.Context, in *pb.GetDeploymentRequest) (*pb.GetDeploymentResponse, error) {
	//TODO
	return &pb.GetDeploymentResponse{Data: nil}, nil
}

func (s *server) DeleteDeployment(ctx context.Context, in *pb.DeleteDeploymentRequest) (*pb.StatusResponse, error) {
	//TODO
	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) ApplyDeployment(ctx context.Context, in *pb.ApplyDeploymentRequest) (*pb.StatusResponse, error) {
	//TODO
	return &pb.StatusResponse{Status: 0}, nil
}

func Run() {
	/**
	**   开启etcd
	**/
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

	/**
	**   创建gRPC服务器,接受来自Kubectl和ApiServer的请求
	**/
	svr := grpc.NewServer()
	// 将实现的接口注册进 gRPC 服务器
	pb.RegisterApiServerKubeletServiceServer(svr, &server{})
	log.Println("Apiserver For Kubelet gRPC Server starts running...")

	pb.RegisterApiServerKubectlServiceServer(svr, &server{})
	log.Println("Apiserver For Kubectl gRpc Server starts running...")

	// 启动 gRPC 服务器
	err = svr.Serve(listen)
	if err != nil {
		log.Fatal(err)
		return
	}
}
