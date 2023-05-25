package app

import (
	"context"
	"encoding/json"
	"fmt"
	"minik8s/configs"
	"minik8s/tools/log"
	"strings"

	"minik8s/entity"
	"minik8s/tools/etcdctl"
	"net"

	"minik8s/pkg/kubelet"
	pb "minik8s/pkg/proto"

	"minik8s/pkg/apiserver"

	//clientv3 "go.etcd.io/etcd/client/v3"

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
	log.Print("[Api Server] Kubelet call sayHello...")

	log.Print(in)
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

// ApplyPod 客户端为Kubectl
func (s *server) ApplyPod(ctx context.Context, in *pb.ApplyPodRequest) (*pb.StatusResponse, error) {

	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	if err != nil {
		log.PrintE("pod unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}

	cli, err := etcdctl.NewClient()
	defer cli.Close()
	if err != nil {
		log.PrintE("etcd client connetc error")
	}
	log.Print("put etcd")
	etcdctl.Put(cli, "Pod/"+pod.Metadata.Name, string(in.Data))

	return apiserver.ApiServerObject().ApplyPod(in)
}

func (s *server) DeletePod(ctx context.Context, in *pb.DeletePodRequest) (*pb.StatusResponse, error) {

	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("connect to etcd error")
	}
	defer cli.Close()
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
		log.PrintE("connect to etcd error")
	}
	defer cli.Close()

	out, err := etcdctl.Get(cli, "Pod/"+string(in.PodName))
	if in.PodName == "" {
		out, err = etcdctl.GetWithPrefix(cli, "Pod/")
	}

	// conver []*mvccpb.KeyValue to []byte
	var data [][]byte
	for _, v := range out.Kvs {
		data = append(data, v.Value)
	}

	if len(out.Kvs) == 0 {
		return &pb.GetPodResponse{PodData: nil}, nil
	} else {
		return &pb.GetPodResponse{PodData: data}, nil
	}
}

func (s *server) GetNode(ctx context.Context, in *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("connect to etcd error")
	}

	defer cli.Close()
	out, _ := etcdctl.Get(cli, "Node/"+string(in.NodeName))
	fmt.Println(out.Kvs)
	if len(out.Kvs) == 0 {
		return &pb.GetNodeResponse{NodeData: nil}, nil
	} else {
		return &pb.GetNodeResponse{NodeData: out.Kvs[0].Value}, nil
	}
}

func (s *server) ApplyJob(ctx context.Context, in *pb.ApplyJobRequest) (*pb.StatusResponse, error) {
	// 解析Job
	job := &entity.Job{}
	err := json.Unmarshal(in.Data, job)
	if err != nil {
		log.PrintE("pod unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}

	// 存入etcd中
	cli, err := etcdctl.NewClient()
	defer cli.Close()
	if err != nil {
		log.PrintE("etcd client connetc error")
	}
	log.Print("put etcd")
	etcdctl.Put(cli, "Job/"+job.Metadata.Name, string(in.Data))

	return apiserver.ApiServerObject().ApplyJob(job)
}

<<<<<<< HEAD
func (s *server)ApplyFunction(ctx context.Context, in *pb.ApplyFunctionRequest) (*pb.StatusResponse, error) {
	// 解析Function
	function := &entity.Function{}
	err := json.Unmarshal(in.Data, function)
	if err != nil {
		log.PrintE("pod unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}
    
	return apiserver.ApiServerObject().ApplyFunction(function)
}

=======
>>>>>>> master
// 客户端为Kubelet
func (s *server) RegisterNode(ctx context.Context, in *pb.RegisterNodeRequest) (*pb.StatusResponse, error) {
	newNode := &entity.Node{}
	newNode.Ip = in.NodeIp
	newNode.Name = in.NodeName
	newNode.KubeletUrl = in.KubeletUrl
	newNode.Status = entity.NodeLive
	apiserver.ApiServerObject().NodeManager.RegiseterNode(newNode)
	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) UpdatePodStatus(ctx context.Context, in *pb.UpdatePodStatusRequest) (*pb.StatusResponse, error) {
	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	if err != nil {
		log.PrintE("pod unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}

	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("etcd client connect error")
	}
	defer cli.Close()

	//检查本地etcd中是否有此Pod，没有说明已经是一个删除的Pod，将其etcd端信息写为succeed,下一次Pod更新通知删除
	response, err := etcdctl.Get(cli, "Pod/"+pod.Metadata.Name)
	if len(response.Kvs) == 0 {
		pod.Status.Phase = entity.Succeed
	}

	log.PrintS("Update Pod Status: put etcd:", string(in.Data))
	etcdctl.Put(cli, "Pod/"+pod.Metadata.Name, string(in.Data))
	//更新deployment replica
	if strings.Contains(pod.Metadata.Name, "deployment") {

		str := pod.Metadata.Name
		index := strings.Index(str, "deployment")
		deploymentName := ""
		if index != -1 {
			deploymentName = str[:index]
		}
		deploymentName = deploymentName + "deployment"
		log.Print("Update deployment Status.Replicas", deploymentName)
		out, err := etcdctl.Get(cli, "Deployment/"+deploymentName)
		if err != nil {
			log.Print("deployment %s not exist", deploymentName)
		}
		deployment := &entity.Deployment{}
		err = json.Unmarshal(out.Kvs[0].Value, deployment)
		deployment.Status.Replicas += 1
		deploymentByte, err := json.Marshal(deployment)
		etcdctl.Put(cli, "Deployment/"+deploymentName, string(deploymentByte))
	}

	return &pb.StatusResponse{Status: 0}, err
}

// GetService Service
func (s *server) GetService(ctx context.Context, in *pb.GetServiceRequest) (*pb.GetServiceResponse, error) {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("connect to etcd error")
	}
	defer cli.Close()

	out, _ := etcdctl.Get(cli, "Service/"+string(in.ServiceName))
	if in.ServiceName == "" {
		out, _ = etcdctl.GetWithPrefix(cli, "Service/")
	}

	// conver []*mvccpb.KeyValue to []byte
	var data [][]byte
	for _, v := range out.Kvs {
		data = append(data, v.Value)
	}

	if len(out.Kvs) == 0 {
		return &pb.GetServiceResponse{Data: nil}, nil
	} else {
		return &pb.GetServiceResponse{Data: data}, nil
	}
}

// GetJob Service
func (s *server) GetJob(ctx context.Context, in *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("connect to etcd error")
	}
	out, _ := etcdctl.Get(cli, "Job/"+string(in.JobName))
	fmt.Println(out.Kvs)
	if len(out.Kvs) == 0 {
		return &pb.GetJobResponse{Data: nil}, nil
	} else {
		return &pb.GetJobResponse{Data: out.Kvs[0].Value}, nil
	}
}

func (s *server) DeleteService(ctx context.Context, in *pb.DeleteServiceRequest) (*pb.StatusResponse, error) {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("connect to etcd error")
	}
	defer cli.Close()
	out, _ := etcdctl.Get(cli, "Service/"+string(in.ServiceName))
	fmt.Println(out.Kvs)
	if len(out.Kvs) == 0 {
		return &pb.StatusResponse{Status: 0}, nil
	}

	err = kubelet.KubeProxyObject().RemoveService(in.ServiceName)
	if err != nil {
		log.PrintE(err)
	}
	// service := &entity.Service{}
	// for _,data := range out.Kvs {
	// 	err := json.Unmarshal(data.Value, service)
	// 	if err != nil {
	// 		log.PrintE("service unmarshal error")
	// 	}

	// 	if service.Metadata.Name == in.ServiceName {
	// 		kubelet.KubeProxyObject().RemoveService(in.ServiceName)
	// 		break;
	// 	}
	// }

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
		log.PrintE("etcd client connetc error")
	}
	defer cli.Close()
	etcdctl.Put(cli, "Service/"+service.Metadata.Name, string(in.Data))

	// 获取符合条件的Pod
	selectedPods := ControllerManager.GetPodsByLabels(&service.Spec.Selector)
	ControllerManager.PrintList(selectedPods)

	// 组装信息
	podNames := make([]string, 0, selectedPods.Len())
	podIps := make([]string, 0, selectedPods.Len())
	for it := selectedPods.Front(); it != nil; it = it.Next() {
		pod := it.Value.(*entity.Pod)
		podNames = append(podNames, pod.Metadata.Name)
		podIps = append(podIps, pod.Status.PodIp)
	}

	if in.Data == nil || podNames == nil || podIps == nil {
		log.PrintE("service data or pod is <nil>")
		log.Print(in.Data)
		log.Print(podNames)
		log.Print(podIps)
	}
	return apiserver.ApiServerObject().CreateService(&pb.ApplyServiceRequest2{
		Data:     in.Data,
		PodNames: podNames,
		PodIps:   podIps,
	})
}

// Deployment
func (s *server) GetDeployment(ctx context.Context, in *pb.GetDeploymentRequest) (*pb.GetDeploymentResponse, error) {
	//TODO
	return &pb.GetDeploymentResponse{Data: nil}, nil
}

func (s *server) DeleteDeployment(ctx context.Context, in *pb.DeleteDeploymentRequest) (*pb.StatusResponse, error) {
	//TODO
	apiserver.ApiServerObject().DeleteDeployment(in)
	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) ApplyDeployment(ctx context.Context, in *pb.ApplyDeploymentRequest) (*pb.StatusResponse, error) {
	//TODO 调用DeploymentController 创建deployment
	apiserver.ApiServerObject().AddDeployment(in)
	return &pb.StatusResponse{Status: 0}, nil
}

// Dns
func (s *server) GetDns(ctx context.Context, in *pb.GetDnsRequest) (*pb.GetDnsResponse, error) {
	// get dns info from etcd
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connetc error")
	}
	defer cli.Close()
	out, _ := etcdctl.Get(cli, "Dns/"+string(in.DnsName))
	// fmt.Println(out.Kvs)
	if len(out.Kvs) == 0 {
		return &pb.GetDnsResponse{Data: nil}, nil
	} else {
		return &pb.GetDnsResponse{Data: out.Kvs[0].Value}, nil
	}

}

func (s *server) DeleteDns(ctx context.Context, in *pb.DeleteDnsRequest) (*pb.StatusResponse, error) {
	return apiserver.ApiServerObject().DeleteDns(in)
}

func (s *server) ApplyDns(ctx context.Context, in *pb.ApplyDnsRequest) (*pb.StatusResponse, error) {
	dns := &entity.Dns{}
	err := json.Unmarshal(in.Data, dns)
	if err != nil {
		log.PrintE(err)
		return &pb.StatusResponse{Status: -1}, err
	}

	// get all services from etcd
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connetc error")
	}
	defer cli.Close()

	// put dns info into etcd
	etcdctl.Put(cli, "Dns/"+dns.Metadata.Name, string(in.Data))

	out, _ := etcdctl.GetWithPrefix(cli, "Service/")
	services := make([]*entity.Service, 0, len(out.Kvs))
	if len(out.Kvs) == 0 {
		log.PrintE("no service found")
		return &pb.StatusResponse{Status: -1}, err
	}
	for _, data := range out.Kvs {
		service := &entity.Service{}
		err := json.Unmarshal(data.Value, service)
		if err != nil {
			log.PrintE("service unmarshal error")
		}
		services = append(services, service)
	}

	// 将dns的serviceName字段换成对应service的clusterIP
	for i, serviceName := range dns.Spec.Paths {
		flag := false
		for _, service := range services {
			if service.Metadata.Name == serviceName.ServiceName {
				dns.Spec.Paths[i].ServiceName = service.Spec.ClusterIP
				flag = true
				break
			}
		}
		if !flag {
			log.PrintE("service not found")
			return &pb.StatusResponse{Status: -1}, err
		}
	}

	data, err := json.Marshal(dns)
	if err != nil {
		log.PrintE(err)
		return &pb.StatusResponse{Status: -1}, err
	}
	return apiserver.ApiServerObject().ApplyDns(&pb.ApplyDnsRequest{
		Data: data,
	})

}

func Run() {
	/**
	**   开启etcd
	**/
	// cli, err := etcdctl.Start(configs.EtcdStartPath)
	// if err != nil {
	// 	return
	// }
	// defer func(cli *clientv3.Client) {
	// 	err := cli.Close()
	// 	if err != nil {
	// 		log.PrintE("etcd close error")
	// 	}
	// }(cli)

	// 注册请求处理接口
	listen, err := net.Listen("tcp", configs.GrpcPort)
	if err != nil {
		log.PrintE(err)
		return
	}

	//启动Pod监控
	//go apiserver.ApiServerObject().BeginMonitorPod()
	//log.PrintS("Apiserver For PodMonitor Server starts running...")
	//
	////启动deployment监控
	//go ControllerManager.BeginMonitorDeployment()
	//log.PrintS("Apiserver For DeploymentMonitor Server starts running...")

	/**
	*  创建Http Trigger
	**/
    

	/**
	**   创建gRPC服务器,接受来自Kubectl和ApiServer的请求
	**/
	svr := grpc.NewServer()
	// 将实现的接口注册进 gRPC 服务器
	pb.RegisterApiServerKubeletServiceServer(svr, &server{})
	log.Print("Apiserver For Kubelet gRPC Server starts running...")

	pb.RegisterApiServerKubectlServiceServer(svr, &server{})
	log.Print("Apiserver For Kubectl gRpc Server starts running...")

	// 启动 gRPC 服务器
	err = svr.Serve(listen)
	if err != nil {
		log.PrintE(err)
		return
	}

}
