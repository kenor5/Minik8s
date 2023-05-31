package apiserver

import (
	"encoding/json"
	"fmt"
	"minik8s/configs"
	"minik8s/pkg/apiserver/ControllerManager"
	"minik8s/pkg/apiserver/ControllerManager/JobController"
	"minik8s/pkg/apiserver/ControllerManager/ScaleController"
	"minik8s/pkg/apiserver/scale"
	"minik8s/tools/etcdctl"
	"strconv"
	"time"

	//clientv3 "go.etcd.io/etcd/client/v3"

	// "minik8s/configs"

	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
	"minik8s/entity"
	Controller "minik8s/pkg/apiserver/ControllerManager"
	"minik8s/pkg/apiserver/ControllerManager/FunctionController"
	"minik8s/pkg/apiserver/ControllerManager/NodeController"
	"minik8s/pkg/apiserver/client"
	"minik8s/pkg/kubelet/container/containerfunc"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
	"os"
	"os/exec"
)

/**************************************************************************
************************    API Server 主结构    ****************************
***************************************************************************/
type ApiServer struct {
	// conn pb.KubeletApiServerServiceClient
	NodeManager     NodeController.NodeController
	FunctionManager functioncontroller.FunctionController
	AM              ScaleController.AutoscalerManager
}

var apiServer *ApiServer

func newApiServer() *ApiServer {
	newServer := &ApiServer{
		NodeManager:     *NodeController.NewNodeController(),
		FunctionManager: *functioncontroller.NewFunctionController(),
		AM: ScaleController.AutoscalerManager{
			MetricsManager: scale.NewMetricsManager(),
			Autoscalers:    map[string]*entity.HorizontalPodAutoscaler{},
		},
	}
	return newServer
}

func ApiServerObject() *ApiServer {
	if apiServer == nil {
		apiServer = newApiServer()
	}
	return apiServer
}

func (master *ApiServer) ApplyPod(in *pb.ApplyPodRequest) (*pb.StatusResponse, error) {
	// 调度(获取conn)
	conn, hostip := master.NodeManager.RoundRobin()
	//更新hostip信息
	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	pod.Status.HostIp = hostip

	// 发送消息给Kubelet
	poddata, _ := json.Marshal(pod)
	etcdctl.EtcdPut("Pod/"+pod.Metadata.Name, string(poddata))
	in.Data = poddata

	err = client.KubeletCreatePod(conn, in)
	if err != nil {
		log.PrintE(err)
		return &pb.StatusResponse{Status: -1}, err
	}

	return &pb.StatusResponse{Status: 0}, err
}

func (master *ApiServer) DeletePod(in *pb.DeletePodRequest) (*pb.StatusResponse, error) {
	//查询Pod对应的Node信息并获取conn
	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	if in.Data == nil || pod.Status.Phase == entity.Succeed {
		return &pb.StatusResponse{Status: 0}, err
	}
	// 根据Pod所在的节点的NodeName获得对应的grpc Conn
	conn := master.NodeManager.GetNodeConnByName(pod.Spec.NodeName)
	if conn == nil {
		conn = master.NodeManager.GetNodeConnByIP(pod.Status.HostIp)
	}
	if conn == nil {
		panic("UnKnown NodeName!\n")
	}
	// 发送消息给Kubelet
	err = client.KubeletDeletePod(conn, in)
	if err != nil {
		log.PrintE(err)

		return &pb.StatusResponse{Status: -1}, err
	}

	return &pb.StatusResponse{Status: 0}, err
}

func (master *ApiServer) CreateService(in *pb.ApplyServiceRequest2) (*pb.StatusResponse, error) {
	LivingNodes := master.NodeManager.GetAllLivingNodes()

	for _, node := range LivingNodes {
		// 发送消息给Kubelet
		conn := master.NodeManager.GetNodeConnByName(node.Name)
		err := client.KubeLetCreateService(conn, in)
		if err != nil {
			log.PrintE(err)
			return &pb.StatusResponse{Status: -1}, err
		}
	}
	return &pb.StatusResponse{Status: 0}, nil

}

func (master *ApiServer) DeleteService(in *pb.DeleteServiceRequest) (*pb.StatusResponse, error) {
	LivingNodes := master.NodeManager.GetAllLivingNodes()

	for _, node := range LivingNodes {
		// 发送消息给Kubelet
		conn := master.NodeManager.GetNodeConnByName(node.Name)
		err := client.KubeLetDeleteService(conn, &pb.DeleteServiceRequest2{
			ServiceName: in.ServiceName,
		})
		if err != nil {
			log.PrintE(err)
			return &pb.StatusResponse{Status: -1}, err
		}

	}
	return &pb.StatusResponse{Status: 0}, nil
}

// AddDeployment 增加新的deployment，先将元数据写入etcd,之后监控机制检查更新pod状态并启动
func (master *ApiServer) AddDeployment(in *pb.ApplyDeploymentRequest) {
	deployment := &entity.Deployment{}
	err := json.Unmarshal(in.Data, deployment)
	if err != nil {
		fmt.Print("[ApiServer]ApplyDeployment Unmarshal error!\n")
		return
	}
	//写入etcd元数据
	_, err = Controller.ApplyDeployment(deployment)
	if err != nil {
		return
	}
	//依次创建deployment 中的pod
	//for _, pod := range podList {
	//	podByte, err := json.Marshal(pod)
	//	if err != nil {
	//		fmt.Println("parse pod error")
	//		return
	//	}
	//	_, err = master.ApplyPod(&pb.ApplyPodRequest{
	//		Data: podByte,
	//	})
	//	if err != nil {
	//		fmt.Printf("create Pod of Deployment error:%s", err)
	//		return
	//	}
	//	//err = client.KubeletCreatePod(apiServer.conn)
	//	//if err != nil {
	//	//	return
	//	//}
	//}
}

// DeleteDeployment  使用Controller删除deployment
func (master *ApiServer) DeleteDeployment(in *pb.DeleteDeploymentRequest) {
	deploymentname := in.DeploymentName
	//从etcd中删除该deployment
	err := Controller.DeleteDeployment(deploymentname)
	if err != nil {
		return
	}
}

func (master *ApiServer) AddHPA(HPAbyte *pb.ApplyHorizontalPodAutoscalerRequest) {
	HPA := &entity.HorizontalPodAutoscaler{}
	err := json.Unmarshal(HPAbyte.Data, HPA)
	if err != nil {
		log.Print("[ApiServer]AddHPA Unmarshal error!\n")
		return
	}
	master.AM.CreateAutoscaler(HPA)
	go master.AM.StartAutoscalerMonitor(HPA)
	//写入etcd元数据
	if err != nil {
		fmt.Println("etcd client connect error")
	}
	HPAData, err := json.Marshal(HPA)
	err = etcdctl.EtcdPut("HPA/"+HPA.Metadata.Name, string(HPAData))
	if err != nil {
		log.PrintE("[ApiServer] Write HPAData fail")
		return
	}
}

func (master *ApiServer) DeleteHPA(HPAbyte *pb.DeleteHorizontalPodAutoscaler) {
	HPAName := HPAbyte.Data
	master.AM.DeleteAutoscaler(HPAName)

	//删除etcd元数据
	etcdctl.EtcdDelete("HPA/" + HPAName)
}

// BeginMonitorPod TODO 根据etcd中存储的信息，监控和更新Pod的状态
func (master *ApiServer) BeginMonitorPod() {

	//检查pending状态的Pod，指定一个Node创建
	ticker := time.NewTicker(configs.MonitorPodTime * time.Second) //每5s检查一次etcd信息
	for range ticker.C {
		master.CheckEtcdAndUpdate()
	}
}
func (master *ApiServer) CheckEtcdAndUpdate() {
	log.Printf("[Apiserver]CheckEtcdAndUpdate begin monitor Pod")
	PodsData, _ := etcdctl.EtcdGetWithPrefix("Pod/")
	//读取所有的Pod信息
	var Pods []entity.Pod
	for _, PodData := range PodsData.Kvs {
		var pod entity.Pod
		err := json.Unmarshal(PodData.Value, &pod)
		if err != nil {
			log.PrintE("[CheckEtcdAndUpdate]GetPods Unmarshal Pod error")
			return
		}
		Pods = append(Pods, pod)
	}
	log.Print("[CheckEtcdAndUpdate]读取所有的Pod信息")
	//根据不同的Pod状态，进行不同处理
	/*
		1. Pending 为创建后还为指定Node实际创建
		2. Failed  为Node上Pod运行已经监控到fail或者Node无法连接
		3. Running 不用处理
		4. Succeed 应该极少出现，因为先删除本地etcd才通知Node删除
	*/
	for _, pod := range Pods {
		podstate := pod.Status.Phase
		switch podstate {
		case entity.Running:
			continue
		case entity.Failed:
			//TODO: 利用hostIP通知检查Node状态,Node存活则会自动更新状态
			//一种简单通用实现：直接发送一次DeletePod 到Node，再为Pod重新指定Node创建
			hostIP := pod.Status.HostIp
			if hostIP == "" {
				//需要分配新的Node
				hostIP = "127.0.0.1"
			}
		case entity.Pending:
			//TODO: 此类为已经写入etcd但还未指定Node创建的Pod，如新的replica
			//根据调度策略选择合适Node分配
			// 组装消息
			log.Print("[CheckEtcdAndUpdate]处理Pod:Pending")

			// 调度(获取conn)
			conn, hostip := master.NodeManager.RoundRobin()
			if conn == nil {
				KubeletUrl := "127.0.0.1:5679"
				conn, err := NodeController.ConnectToKubelet("127.0.0.1:5679")
				if err != nil {
					panic("fail to connect kubelet: " + KubeletUrl)
				}
				pod.Status.HostIp = "127.0.0.1"
				podByte, err := json.Marshal(pod)
				if err != nil {
					fmt.Println("parse pod error")
					continue
				}
				in := &pb.ApplyPodRequest{
					Data: podByte,
				}
				err = etcdctl.EtcdPut("Pod/"+pod.Metadata.Name, string(podByte))
				if err != nil {
					log.PrintE("Put Pod error")
				}
				err = client.KubeletCreatePod(conn, in)
			} else {
				// 发送消息给Kubelet
				if hostip == "" {
					log.PrintE("Host ip is nil")
					panic(hostip)
					//hostip = "127.0.0.1"
				}
				pod.Status.HostIp = hostip
				podByte, err := json.Marshal(pod)
				if err != nil {
					fmt.Println("parse pod error")
					continue
				}
				in := &pb.ApplyPodRequest{
					Data: podByte,
				}
				err = etcdctl.EtcdPut("Pod/"+pod.Metadata.Name, string(podByte))
				if err != nil {
					log.PrintE("Put Pod error")
				}
				err = client.KubeletCreatePod(conn, in)
			}

		case entity.Succeed:
			//TODO: 通知Node删除Pod后，更新本地信息
			if pod.Status.HostIp == "" {
				//已经完全删除退出的Pod
				continue
			}
			log.Print("[CheckEtcdAndUpdate]处理Pod:Succeed")
			conn := master.NodeManager.GetNodeConnByIP(pod.Status.HostIp)
			pod.Status.HostIp = ""
			pod.Status.PodIp = ""
			podByte, _ := json.Marshal(pod)
			err := client.KubeletDeletePod(conn, &pb.DeletePodRequest{
				Data: podByte,
			})
			if err != nil {
				log.PrintfE("[CheckEtcdAndUpdate]Delete Pod %s error", pod.Metadata.Name)
			}

			//err = etcdctl.EtcdDelete("Pod/" + pod.Metadata.Name)
			//if err != nil {
			//	log.PrintfE("[CheckEtcdAndUpdate]Delete Pod %s error", pod.Metadata.Name)
			//	continue
			//}
		}
	}
}

func (master *ApiServer) ApplyDns(in *pb.ApplyDnsRequest) (*pb.StatusResponse, error) {
	LivingNodes := master.NodeManager.GetAllLivingNodes()

	for _, node := range LivingNodes {
		// 发送消息给Kubelet
		conn := master.NodeManager.GetNodeConnByName(node.Name)
		err := client.KubeLetCreateDns(conn, in)
		if err != nil {
			log.PrintE(err)
			return &pb.StatusResponse{Status: -1}, err
		}

	}
	return &pb.StatusResponse{Status: 0}, nil
}

func (master *ApiServer) DeleteDns(in *pb.DeleteDnsRequest) (*pb.StatusResponse, error) {
	LivingNodes := master.NodeManager.GetAllLivingNodes()

	for _, node := range LivingNodes {
		// 发送消息给Kubelet
		conn := master.NodeManager.GetNodeConnByName(node.Name)
		err := client.KubeLetDeleteDns(conn, in)
		if err != nil {
			log.PrintE(err)
			return &pb.StatusResponse{Status: -1}, err
		}

	}
	return &pb.StatusResponse{Status: 0}, nil
}
func (master *ApiServer) ApplyJob(job *entity.Job) (*pb.StatusResponse, error) {
	pod := &entity.Pod{
		Kind: "pod",
		Metadata: entity.ObjectMeta{
			Name: job.Metadata.Name + "-ServerPod",
			Labels: map[string]string{
				"app": "Job",
			},
		},
		Spec: entity.PodSpec{
			Containers: []entity.Container{
				{
					Name:  "slurm-server",
					Image: "luoshicai/slurm-server:latest",
					VolumeMounts: []entity.VolumeMount{
						{
							Name:      "volume1",
							MountPath: "/tryData",
						},
					},
				},
			},
			Volumes: []entity.Volume{
				{
					Name:     "volume1",
					HostPath: "/root/go/src/minik8s/tools/cuda/" + job.Metadata.Name,
				},
			},
		},
	}

	// 调度(获取conn)
	conn, hostip := master.NodeManager.RoundRobin()
	//更新hostip信息
	pod.Status.HostIp = hostip
	// 组装消息
	podByte, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("parse pod error")
		return &pb.StatusResponse{Status: -1}, err
	}
	in := &pb.ApplyPodRequest{
		Data: podByte,
	}
	// 发送消息给Kubelet
	err = client.KubeletCreatePod(conn, in)
	if err != nil {
		log.PrintE(err)
		return &pb.StatusResponse{Status: -1}, err
	}

	go JobController.SbatchAndQuery(job.Metadata.Name, conn)

	return &pb.StatusResponse{Status: 0}, err
}

func (master *ApiServer) ApplyFunction(function *entity.Function) (*pb.StatusResponse, error) {
	imageName := "luoshicai/" + function.Metadata.Name + ":latest"
	exist, _ := containerfunc.ImageExist(imageName)
	if exist == false {
		log.Print("function image doesn't exist, create function image...")
		// 生成Dockerfile
		// 命令和参数
		cmd := exec.Command("./scripts/gen_function_dockerfile.sh", function.Metadata.Name)
		// 设置工作目录
		cmd.Dir = "./"
		// 执行命令并捕获输出
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.PrintE("命令执行失败：%v\n%s", err, output)
		}

		// 打印输出结果
		log.Print(string(output))

		// 生成镜像
		dockerfilePath := "./tools/serverless/" + function.Metadata.Name

		cmd = exec.Command("docker", "build", "-t", imageName, dockerfilePath)

		// 设置命令输出到标准输出和标准错误
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// 执行命令
		err = cmd.Run()
		if err != nil {
			log.PrintE("执行命令失败：%v", err)
		}

		log.Print("镜像构建成功：%s\n", imageName)
	}

	// 存入etcd
	PodTemplate := &entity.Pod{
		Kind: "pod",
		Metadata: entity.ObjectMeta{
			Name: function.Metadata.Name + "-pod",
			Labels: map[string]string{
				"app": function.Metadata.Name,
			},
		},
		Spec: entity.PodSpec{
			Containers: []entity.Container{
				{
					Name:  "serverless-server",
					Image: imageName,
				},
			},
		},
	}
	function.FunctionStatus.AccessTimes = 0
	function.FunctionStatus.Status = entity.Running
	function.FunctionStatus.PodTemplate = *PodTemplate
	functioncontroller.SetFunction(function)

	// 加入路由
	master.AddRouter(function.Metadata.Name)

	return &pb.StatusResponse{Status: 0}, nil
}

func (master *ApiServer) ApplyWorkflow(workflow *entity.Workflow) (*pb.StatusResponse, error) {
	// TODO: 判断etcd中是否有Workflow对应的function

	// 存入etcd
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connetc error")
	}
	defer cli.Close()
	workflowByte, err := json.Marshal(workflow)
	etcdctl.Put(cli, "Workflow/"+workflow.Name, string(workflowByte))

	// 加入路由
	master.AddWorkflowRouter(workflow.Name)

	return &pb.StatusResponse{Status: 0}, nil
}

func (master *ApiServer) DeleteFunction(functionName string) (*pb.StatusResponse, error) {
    // 从etcd中查找function
	function, _ := functioncontroller.GetFunction(functionName)
	
	// 将该function从etcd中删除
    err := functioncontroller.DelFunction(functionName)
	if err != nil {
		log.PrintS("Delete function err!")
		return &pb.StatusResponse{Status: -1}, err
	}

	// 删除所有属于这个function的Pod
	functionPods := function.FunctionStatus.FunctionPods
    for _, functionPod := range functionPods {
		pod, _ := Controller.GetPodByName(functionPod.PodName)
		podByte, err := json.Marshal(pod)
		if err != nil {
			log.PrintE("parse pod error")
			return &pb.StatusResponse{Status: -1}, err
		}
		in := &pb.DeletePodRequest{
			Data: podByte,
		}		
		_, err = master.DeletePod(in)
		if err != nil {
			log.PrintE("delete pod error")
			return &pb.StatusResponse{Status: -1}, err
		}
	}

	log.PrintS("delete function: ", functionName)

	return &pb.StatusResponse{Status: 0}, nil
}

func (master *ApiServer) UpdateFunction(functionName string) (*pb.StatusResponse, error) {
	imageName := "luoshicai/" + functionName + ":latest"
    // 从etcd中查找function
	function, _ := functioncontroller.GetFunction(functionName)
	
	// 将该function从etcd中删除
    err := functioncontroller.DelFunction(functionName)
	if err != nil {
		log.PrintS("Delete function err!")
		return &pb.StatusResponse{Status: -1}, err
	}

	// 删除原有镜像
	containerfunc.DeleteImage(imageName)

	// 创建新镜像
	log.Print("begin update function image...")
	// 生成Dockerfile
	// 命令和参数
	cmd := exec.Command("./scripts/gen_function_dockerfile.sh", function.Metadata.Name)
		// 设置工作目录
	cmd.Dir = "./"
	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.PrintE("命令执行失败：%v\n%s", err, output)
	} 

	// 打印输出结果
	log.Print(string(output))

	// 生成镜像
	dockerfilePath := "./tools/serverless/" + function.Metadata.Name

	cmd = exec.Command("docker", "build", "-t", imageName, dockerfilePath)

	// 设置命令输出到标准输出和标准错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	err = cmd.Run()
	if err != nil {
		log.PrintE("执行命令失败：%v", err)
	}

	log.Print("镜像更新成功：%s\n", imageName)	

	// 更新Pod:策略为先删除所有的Pod，再创建同样数量的新Pod
	podNum := len(function.FunctionStatus.FunctionPods)
	for _, functionPod := range function.FunctionStatus.FunctionPods {
		pod, _ := Controller.GetPodByName(functionPod.PodName)
		podByte, err := json.Marshal(pod)
		if err != nil {
			log.PrintE("parse pod error")
			return &pb.StatusResponse{Status: -1}, err
		}
		in := &pb.DeletePodRequest{
			Data: podByte,
		}		
		_, err = master.DeletePod(in)
		if err != nil {
			log.PrintE("delete pod error")
			return &pb.StatusResponse{Status: -1}, err
		}
	}

	newPod := function.FunctionStatus.PodTemplate
	newPodName := newPod.Metadata.Name
	for i := 0; i < podNum; i++ {
		newPod.Metadata.Name = newPodName + "-" + strconv.Itoa(i)
		// 组装消息
		podByte, err := json.Marshal(newPod)
		if err != nil {
		log.PrintE("parse pod error")
		}
		in := &pb.ApplyPodRequest{
			Data: podByte,
		}			
		_, err = master.ApplyPod(in)
		if err != nil {
			log.PrintE("Apply pod error")
		}        
	}
	
	// 然后更新function
	PodsList := ControllerManager.GetPodsByLabels(&newPod.Metadata.Labels)
	targetFunctionPods := []entity.FunctionPod{}
	// 遍历列表
	for it := PodsList.Front(); it != nil; it = it.Next() {
		element := it.Value.(*entity.Pod)
		targetFunctionPods = append(targetFunctionPods, entity.FunctionPod{
			PodName: element.Metadata.Name,
			PodIp: element.Status.PodIp,
		})
	}

	// 将function重新放入etcd中
    functioncontroller.SetFunction(function)

	return &pb.StatusResponse{Status: 0}, nil
}
