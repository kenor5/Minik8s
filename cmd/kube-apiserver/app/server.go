package app

import (
	"context"
	"encoding/json"
	"fmt"
	"minik8s/configs"
	"minik8s/pkg/apiserver/scale"
	"minik8s/tools/log"
	"strings"

	"minik8s/entity"
	"minik8s/tools/etcdctl"
	"net"

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
		log.PrintE("pod unmarshel err %v", err)
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

// GetPod: get pods后不跟PodName返回所有的Pod
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
	if in.NodeName == "" {
		out, _ = etcdctl.GetWithPrefix(cli, "Node/")
	}
	// fmt.Println(out.Kvs)

	var data [][]byte
	for _, v := range out.Kvs {
		data = append(data, v.Value)
	}
	if len(out.Kvs) == 0 {
		return &pb.GetNodeResponse{NodeData: nil}, nil
	} else {
		return &pb.GetNodeResponse{NodeData: data}, nil
	}
}

func (s *server)AddNode(ctx context.Context, in *pb.AddNodeRequest) (*pb.StatusResponse, error) {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("connect to etcd error")
	}

	defer cli.Close()
	out, _ := etcdctl.Get(cli, "Node/"+string(in.NodeName))
	if in.NodeName == "" {
		out, _ = etcdctl.GetWithPrefix(cli, "Node/")
	}

	for _, v := range out.Kvs {
		node := &entity.Node{}
		err := json.Unmarshal(v.Value, node)
		if err != nil {
			log.PrintE("podNew unmarshel err")
			return &pb.StatusResponse{Status: -1}, err
		}

		node.Status = entity.NodeLive
        
		nodeByte, _ := json.Marshal(node)
		etcdctl.Put(cli, "Node/"+node.Name, string(nodeByte))

	}  

	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) DeleteNode(ctx context.Context, in *pb.DeleteNodeRequest) (*pb.StatusResponse, error) {

	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("connect to etcd error")
	}
	defer cli.Close()
	out, err := etcdctl.Get(cli, "Node/"+string(in.NodeName))

	for _, v := range out.Kvs {
		node := &entity.Node{}
		err := json.Unmarshal(v.Value, node)
		if err != nil {
			log.PrintE("podNew unmarshel err")
			return &pb.StatusResponse{Status: -1}, err
		}

		node.Status = entity.NodePending
        
		nodeByte, _ := json.Marshal(node)
		etcdctl.Put(cli, "Node/"+node.Name, string(nodeByte))
	}  

	return &pb.StatusResponse{Status: 0}, nil
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

func (s *server) ApplyFunction(ctx context.Context, in *pb.ApplyFunctionRequest) (*pb.StatusResponse, error) {
	// 解析Function
	function := &entity.Function{}
	err := json.Unmarshal(in.Data, function)
	if err != nil {
		log.PrintE("function unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}

	return apiserver.ApiServerObject().ApplyFunction(function)
}

func (s *server) GetFunction(ctx context.Context, in *pb.GetFunctionRequest) (*pb.GetFunctionResponse, error) {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("connect to etcd error")
	}

	defer cli.Close()
	out, _ := etcdctl.Get(cli, "Function/"+string(in.FunctionName))
	if in.FunctionName == "" {
		out, _ = etcdctl.GetWithPrefix(cli, "Function/")
	}
	// fmt.Println(out.Kvs)

	var data [][]byte
	for _, v := range out.Kvs {
		data = append(data, v.Value)
	}
	if len(out.Kvs) == 0 {
		return &pb.GetFunctionResponse{Data: nil}, nil
	} else {
		return &pb.GetFunctionResponse{Data: data}, nil
	}
}

func (s *server)UpdateFunction(ctx context.Context, in *pb.UpdateFunctionRequest) (*pb.StatusResponse, error) {
	// 获取FunctionName
	functionName := string(in.FunctionName)
    
	return apiserver.ApiServerObject().UpdateFunction(functionName)
}

func (s *server) ApplyWorkflow(ctx context.Context, in *pb.ApplyWorkflowRequest) (*pb.StatusResponse, error) {
	// 解析Wokflow
	workflow := &entity.Workflow{}
	err := json.Unmarshal(in.Data, workflow)
	if err != nil {
		log.PrintE("workflow unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}

	return apiserver.ApiServerObject().ApplyWorkflow(workflow)
}

func (s *server)DeleteFunction(ctx context.Context, in *pb.DeleteFunctionRequest) (*pb.StatusResponse, error) {
	// 获取FunctionName
	functionName := string(in.FunctionName)
    
	return apiserver.ApiServerObject().DeleteFunction(functionName)        
}


// 客户端为Kubelet
func (s *server) RegisterNode(ctx context.Context, in *pb.RegisterNodeRequest) (*pb.RegisterNodeResponse, error) {
	newNode := &entity.Node{}
	err := json.Unmarshal(in.NodeData, newNode)
	if err != nil {
		log.PrintE("workflow unmarshel err")
	}
	// newNode.Ip = in.NodeIp
	// newNode.Name = in.NodeName
	// newNode.KubeletUrl = in.KubeletUrl
	newNode.Status = entity.NodePending
	apiserver.ApiServerObject().NodeManager.RegiseterNode(newNode)
	podsByte := getPodbyHostIP(newNode.Ip)
	return &pb.RegisterNodeResponse{PodData: podsByte}, nil
}

func getPodbyHostIP(hostIP string) [][]byte {
	//var pods []entity.Pod
	bytes := [][]byte{}
	response, _ := etcdctl.EtcdGetWithPrefix("Pod/")
	for _, poddata := range response.Kvs {
		podNew := entity.Pod{}
		err := json.Unmarshal(poddata.Value, &podNew)
		if err != nil {
			log.PrintE("podNew unmarshel err")
			return nil
		}
		if podNew.Status.HostIp == hostIP {
			//pods = append(pods, podNew)
			podByte, _ := json.Marshal(podNew)
			bytes = append(bytes, podByte)
		}
	}

	return bytes
}

func (s *server) UpdatePodStatus(ctx context.Context, in *pb.UpdatePodStatusRequest) (*pb.StatusResponse, error) {
	podNew := &entity.Pod{}
	err := json.Unmarshal(in.Data, podNew)
	if err != nil {
		log.PrintE("podNew unmarshel err")
		return &pb.StatusResponse{Status: -1}, err
	}

	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("etcd client connect error")
	}
	defer cli.Close()

	//检查本地etcd中是否有此Pod，没有说明已经是一个删除的Pod，将其etcd端信息写为succeed,下一次Pod更新通知删除
	response, err := etcdctl.Get(cli, "Pod/"+podNew.Metadata.Name)
	if len(response.Kvs) == 0 {
		log.Print("[UpdatePodStatus]更新一个没有的Pod")
		podNew.Status.Phase = entity.Succeed
		podNew.Status.HostIp = ""
	}
	podOld := &entity.Pod{}
	//仅更新Phase和PodIP
	err = json.Unmarshal(response.Kvs[0].Value, podOld)
	podOld.Status.Phase = podNew.Status.Phase
	podOld.Status.PodIp = podNew.Status.PodIp
	podOld.Status.StartTime = podNew.Status.StartTime
	if podNew.Status.Phase == entity.Succeed {
		podOld.Status.HostIp = ""
	}
	podData, _ := json.Marshal(podOld)
	log.PrintS("Update Pod Status: put etcd:", string(podData))
	//podNew:=&entity.Pod{}

	etcdctl.Put(cli, "Pod/"+podNew.Metadata.Name, string(podData))
	//更新deployment replica
	if strings.Contains(podNew.Metadata.Name, "deployment") {

		str := podNew.Metadata.Name
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
			return nil, err
		}
		deployment := &entity.Deployment{}
		err = json.Unmarshal(out.Kvs[0].Value, deployment)
		if podNew.Status.Phase == entity.Running {
			deployment.Status.Replicas += 1
		} else if podNew.Status.Phase == entity.Succeed {
			deployment.Status.Replicas -= 1
		}
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
	if in.JobName == "" {
		out, _ = etcdctl.GetWithPrefix(cli, "Job/")
	}

	// conver []*mvccpb.KeyValue to []byte
	var data [][]byte
	for _, v := range out.Kvs {
		data = append(data, v.Value)
	}
	if len(out.Kvs) == 0 {
		return &pb.GetJobResponse{Data: nil}, nil
	} else {
		return &pb.GetJobResponse{Data: data}, nil
	}
}

func (s *server) DeleteService(ctx context.Context, in *pb.DeleteServiceRequest) (*pb.StatusResponse, error) {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("connect to etcd error")
	}
	defer cli.Close()
	out, _ := etcdctl.Get(cli, "Service/"+string(in.ServiceName))
	if len(out.Kvs) == 0 {
		log.PrintE("service %s not exist", in.ServiceName)
		return &pb.StatusResponse{Status: 0}, nil
	}

	return apiserver.ApiServerObject().DeleteService(in)
	// err = kubelet.KubeProxyObject().RemoveService(in.ServiceName)
	// if err != nil {
	// 	log.PrintE(err)
	// }
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

	// return &pb.StatusResponse{Status: 0}, nil
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

	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("connect to etcd error")
	}
	defer cli.Close()
	out, _ := etcdctl.Get(cli, "Deployment/"+string(in.DeploymentName))
	if in.DeploymentName == "" {
		out, _ = etcdctl.GetWithPrefix(cli, "Deployment/")
	}
	var data [][]byte
	for _, v := range out.Kvs {
		data = append(data, v.Value)
	}
	if len(out.Kvs) != 0 {
		return &pb.GetDeploymentResponse{Data: data}, nil
	} else {
		return &pb.GetDeploymentResponse{Data: nil}, nil
	}
}

func (s *server) DeleteDeployment(ctx context.Context, in *pb.DeleteDeploymentRequest) (*pb.StatusResponse, error) {
	apiserver.ApiServerObject().DeleteDeployment(in)
	return &pb.StatusResponse{Status: 0}, nil
}

func (s *server) ApplyDeployment(ctx context.Context, in *pb.ApplyDeploymentRequest) (*pb.StatusResponse, error) {
	apiserver.ApiServerObject().AddDeployment(in)
	return &pb.StatusResponse{Status: 0}, nil
}

// HPA
func (s *server) ApplyHPA(ctx context.Context, in *pb.ApplyHorizontalPodAutoscalerRequest) (*pb.StatusResponse, error) {
	apiserver.ApiServerObject().AddHPA(in)
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
	if in.DnsName == "" {
		out, _ = etcdctl.GetWithPrefix(cli, "Dns/")
	}
	// fmt.Println(out.Kvs)

	var data [][]byte
	for _, v := range out.Kvs {
		data = append(data, v.Value)
	}

	if len(out.Kvs) == 0 {
		return &pb.GetDnsResponse{Data: nil}, nil
	} else {
		return &pb.GetDnsResponse{Data: data}, nil
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
    // 初始化Node
	err := apiserver.ApiServerObject().RestartApiserver()
    if (err != nil) {
		log.PrintE("fail to RestartApiserver!")
	}
    log.PrintS("restart apiserver successfully!")


	// 注册请求处理接口
	listen, err := net.Listen("tcp", configs.GrpcPort)
	if err != nil {
		log.PrintE(err)
		return
	}
	//启动Promtheus服务
	err = scale.StartPrometheusServer()
	if err != nil {
		log.PrintE("[Apiserver] Start PrometheusServer error...")
	}
	//启动Pod监控
	go apiserver.ApiServerObject().BeginMonitorPod()
	log.PrintS("Apiserver For PodMonitor Server starts running...")

	//启动deployment监控
	go ControllerManager.BeginMonitorDeployment()
	log.PrintS("Apiserver For DeploymentMonitor Server starts running...")

	// 启动Node监控
	go apiserver.ApiServerObject().MonitorNode()

	/**
	*  Serverless: 创建Http Trigger
	**/
	go apiserver.ApiServerObject().FunctionManager.FunctionServer()

	/**
	*  Serverless: 启动对Function的监控
	**/
	go apiserver.ApiServerObject().MonitorFunction()

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
