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

// 客户端为Kubectl
func (s *server) ApplyPod(ctx context.Context, in *pb.ApplyPodRequest) (*pb.StatusResponse, error) {

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

	return apiserver.ApiServerObject().ApplyPod(in)
}

func (s *server) DeletePod(ctx context.Context, in *pb.DeletePodRequest) (*pb.StatusResponse, error) {

	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("connect to etcd error")
	}
	out, err:= etcdctl.Get(cli, "Pod/"+string(in.Data))
	// fmt.Println(out.Kvs[0].Value)
	// pod := &entity.Pod{}
	// err = json.Unmarshal(out.Kvs[0].Value, pod)
	// if err != nil {
	// 	fmt.Println("pod unmarshal error")
	// }
	// fmt.Println("get etcd", pod)

	return apiserver.ApiServerObject().DeletePod(&pb.DeletePodRequest{
		Data : out.Kvs[0].Value,
	})
}

// TODO: get pods后不跟PodName返回所有的Pod
func (s *server) GetPod(ctx context.Context, in *pb.GetPodRequest) (*pb.GetPodResponse, error) {
	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("connect to etcd error")
	}
	out, err:= etcdctl.Get(cli, "Pod/"+string(in.PodName))
	//fmt.Println(out.Kvs[0].Value)
	// pod := &entity.Pod{}
	// err = json.Unmarshal(out.Kvs[0].Value, pod)
	// if err != nil {
	// 	fmt.Println("pod unmarshal error")
	// }
	// fmt.Println("get etcd", pod)	
	
	return &pb.GetPodResponse{PodData: out.Kvs[0].Value}, nil
}

// 客户端为Kubelet
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
