package apiserver

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/configs"
	"minik8s/tools/etcdctl"
	"time"

	// "minik8s/configs"

	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"

	"minik8s/entity"
	Controller "minik8s/pkg/apiserver/ControllerManager"
	"minik8s/pkg/apiserver/ControllerManager/NodeController"
	"minik8s/pkg/apiserver/client"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
)

/**************************************************************************
************************    API Server 主结构    ****************************
***************************************************************************/
type ApiServer struct {
	// conn pb.KubeletApiServerServiceClient
	NodeManager NodeController.NodeController
}

var apiServer *ApiServer

func newApiServer() *ApiServer {
	newServer := &ApiServer{
		NodeManager: *NodeController.NewNodeController(),
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
	conn := master.NodeManager.RoundRobin()
	// 发送消息给Kubelet
	err := client.KubeletCreatePod(conn, in)
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

// AddDeployment 增加新的deployment，先将元数据写入etcd,之后监控机制检查更新pod状态并启动
func (master *ApiServer) AddDeployment(in *pb.ApplyDeploymentRequest) {
	deployment := &entity.Deployment{}
	err := json.Unmarshal(in.Data, deployment)
	if err != nil {
		fmt.Print("[ApiServer]ApplyDeployment Unmarshal error!\n")
		return
	}
	//写入etcd元数据
	podList, err := Controller.ApplyDeployment(deployment)
	if err != nil {
		return
	}
	//依次创建deployment 中的pod
	for _, pod := range podList {
		podByte, err := json.Marshal(pod)
		if err != nil {
			fmt.Println("parse pod error")
			return
		}
		_, err = master.ApplyPod(&pb.ApplyPodRequest{
			Data: podByte,
		})
		if err != nil {
			fmt.Printf("create Pod of Deployment error:%s", err)
			return
		}
		//err = client.KubeletCreatePod(apiServer.conn)
		//if err != nil {
		//	return
		//}
	}
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

func (master *ApiServer) ApplyHPA(HPAbyte *pb.ApplyHorizontalPodAutoscalerRequest) {
	HPA := &entity.Deployment{}
	err := json.Unmarshal(HPAbyte.Data, HPA)
	if err != nil {
		fmt.Print("[ApiServer]ApplyHPA Unmarshal error!\n")
		return
	}
	//写入etcd元数据
	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("etcd client connect error")
	}
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			fmt.Print("close etcdClient error!")
		}
	}(cli)
	HPAData, err := json.Marshal(HPA)
	err = etcdctl.Put(cli, "HPA/"+HPA.Metadata.Name, string(HPAData))
	if err != nil {
		log.PrintE("[ApiServer] Write HPAData fail")
		return
	}
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
			//TODO: 利用hostIP通知Node检查状态或者重新创建
			//一种简单通用实现：直接发送一次DeletePod 到Node，再为Pod重新指定Node创建
			hostIP := pod.Status.HostIp
			if hostIP == "" {
				//需要分配新的Node
				hostIP = "127.0.0.1"
			}

		case entity.Pending:
		//TODO: 此类为已经写入etcd但还未指定Node创建的Pod，如新的replica
		//根据调度策略选择合适Node分配

		case entity.Succeed:
			//TODO: 通知Node删除Pod后，删除本地信息
			err := etcdctl.EtcdDelete("Pod/" + pod.Metadata.Name)
			if err != nil {
				log.PrintfE("[CheckEtcdAndUpdate]Delete Pod %s error", pod.Metadata.Name)
				continue
			}
		}
	}

}
