package apiserver

import (
	"encoding/json"
	"fmt"
	// "minik8s/configs"

	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
	"os"
	"os/exec"
	"minik8s/entity"
	"minik8s/pkg/apiserver/ControllerManager/NodeController"
	Controller "minik8s/pkg/apiserver/ControllerManager"
	"minik8s/pkg/apiserver/client"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
	"minik8s/pkg/apiserver/ControllerManager/JobController"
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
	if (in.Data == nil || pod.Status.Phase == entity.Succeed) {
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

// AddDeployment TODO 修改Controller的使用逻辑？？？
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
	Controller.DeleteDeployment(deploymentname)
}

func (master *ApiServer) ApplyJob(job *entity.Job) (*pb.StatusResponse, error) {
    pod := &entity.Pod{
		Kind : "pod",
		Metadata : entity.ObjectMeta{
			Name : job.Metadata.Name + "-ServerPod", 
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
			Volumes : []entity.Volume{
			    {
				    Name : "volume1",
				    HostPath: "/home/luoshicai/go/src/minik8s/tools/cuda/"+job.Metadata.Name,
			    },
			},
		},
	}

	// 组装消息
	podByte, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("parse pod error")
		return &pb.StatusResponse{Status: -1}, err
	}
	in := &pb.ApplyPodRequest{
		Data : podByte,
	}
	// 调度(获取conn)
	conn := master.NodeManager.RoundRobin()
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
	imageName := "luoshicai/" + function.Metadata.Name
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

	
	return &pb.StatusResponse{Status: 0}, nil
}