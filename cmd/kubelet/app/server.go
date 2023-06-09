package app

import (
	"context"
	"encoding/json"
	"minik8s/configs"
	"minik8s/entity"
	"minik8s/pkg/kubelet"
	res "minik8s/pkg/kubelet/resourse"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
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
	log.Print("[Kubelet] Api Server call sayHello...")

	log.Print(in)
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

func (s *server) CreatePod(ctx context.Context, in *pb.ApplyPodRequest) (*pb.StatusResponse, error) {
	log.Print("[Kubelet] Api Server call Create Pod...")
	log.Print(in)

	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	if err != nil {
		log.PrintE("pod unmarshal error")
		return &pb.StatusResponse{Status: -1}, err
	}

	// 调用kubelet对象实现创建Pod的真正逻辑
	err = kubelet.KubeletObject().CreatePod(pod)

	if err != nil {
		log.PrintE("create pod error")
		return &pb.StatusResponse{Status: -1}, err
	}
	log.PrintS("[Kubelet] Create Pod Success")
	return &pb.StatusResponse{Status: 0}, err
}

func (s *server) DeletePod(ctx context.Context, in *pb.DeletePodRequest) (*pb.StatusResponse, error) {
	log.Print("[Kubelet] Api Server call delete Pod...")
	log.Print(in)

	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	if err != nil {
		log.PrintE("pod unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}

	// 调用kubelet对象实现删除 Pod的真正逻辑
	err = kubelet.KubeletObject().DeletePod(pod)

	if err != nil {
		log.PrintE("delete pod error")
		return &pb.StatusResponse{Status: -1}, err
	}

	log.PrintS("[Kubelet] delete Pod Success")
	return &pb.StatusResponse{Status: 0}, err
}

func (s *server) CreateService(ctx context.Context, in *pb.ApplyServiceRequest2) (*pb.StatusResponse, error) {
	service := &entity.Service{}
	err := json.Unmarshal(in.Data, service)
	if err != nil {
		log.PrintE("create service: umarshal service failed")
		return &pb.StatusResponse{Status: -1}, err
	}

	err = kubelet.KubeProxyObject().NewService(service, in.PodNames, in.PodIps)
	if err != nil {
		log.PrintE("kubeproxy create service failed")
		return &pb.StatusResponse{Status: -1}, err
	}

	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) DeleteService(ctx context.Context, in *pb.DeleteServiceRequest2) (*pb.StatusResponse, error) {
	err := kubelet.KubeProxyObject().RemoveService(in.ServiceName)
	if err != nil {
		log.PrintE("kubelet delete service failed")
		return &pb.StatusResponse{Status: -1}, nil
	}

	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) AddPod2Service (ctx context.Context, in *pb.AddPod2ServiceRequest) (*pb.StatusResponse, error) {
	err := kubelet.KubeProxyObject().AddPod2Service(in.ServiceName,  in.PodName, in.PodIp,uint32(in.TargetPort))
	if err != nil {
		log.PrintE(err)
		log.PrintE("kubelet add pod 2 service failed")
		return &pb.StatusResponse{Status: -1}, nil
	}

	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) RemovePodFromService (ctx context.Context, in *pb.RemovePodFromServiceRequest) (*pb.StatusResponse, error) {
	err := kubelet.KubeProxyObject().RemovePodFromService(in.ServiceName, in.PodName)
	if err != nil {
		log.PrintE(err)
		log.PrintE("kubelet remove pod from service failed")
		return &pb.StatusResponse{Status: -1}, nil
	}

	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) CreateDns(ctx context.Context, in *pb.ApplyDnsRequest) (*pb.StatusResponse, error) {
	dns := &entity.Dns{}
	err := json.Unmarshal(in.Data, dns)
	if err != nil {
		log.PrintE("create dns: umarshal dns failed")
		return &pb.StatusResponse{Status: -1}, err
	}

	err = kubelet.KubeletObject().ApplyDns(dns)
	if err != nil {
		log.PrintE("kubelet create dns failed")
		return &pb.StatusResponse{Status: -1}, nil
	}

	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) DeleteDns(ctx context.Context, in *pb.DeleteDnsRequest) (*pb.StatusResponse, error) {
	err := kubelet.KubeletObject().DeleteDns(in.DnsName)
	if err != nil {
		log.PrintE("kubelet delete dns failed")
		return &pb.StatusResponse{Status: -1}, nil
	}

	return &pb.StatusResponse{Status: 0}, nil
}

/*********************************************************
**********************  Kubelet主程序   *************************
**********************************************************/
func Run() {
	log.Print("[kubelet] running...")

	_ = kubelet.KubeProxyObject()
	/**
	 *    Kubelet启动时向APIServer注册
	 **/
	kubelet.KubeletObject().RegisterNode()
	log.PrintS("[kubelet] has registered to apiserver...")

	/**
	 *    Kubelet启动自己的服务端，接受来自ApiServer的消息
	**/
	listen, err := net.Listen("tcp", configs.KubeletGrpcPort)
	if err != nil {
		log.PrintE(err)
	}
	//启动cAdvisor
	err = res.StartcAdvisor()
	if err != nil {
		log.PrintE("[kubelet] Start cAdvisor error...")
	}

	// 开启Pod状态监控和更新
	go kubelet.KubeletObject().BeginMonitor()
	log.PrintS("[kubelet]Begin Monitor Deployment")

	// 创建gRPC服务器
	svr := grpc.NewServer()
	// 将实现的接口注册进 gRPC 服务器
	pb.RegisterKubeletApiServerServiceServer(svr, &server{})

	// 启动 gRPC 服务器
	err = svr.Serve(listen)
	if err != nil {
		log.PrintE(err)
		return
	}
}
