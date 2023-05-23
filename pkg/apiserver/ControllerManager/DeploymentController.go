package ControllerManager

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/entity"
	"minik8s/tools/etcdctl"
	HASH "minik8s/tools/hash"
	log "minik8s/tools/log"
	UUID "minik8s/tools/uuid"
	"strconv"
	"time"
)

// ApplyDeployment 将对应的pod抽象和deployment抽象写入etcd
func ApplyDeployment(deployment *entity.Deployment) ([]*entity.Pod, error) {
	Replicas := 0
	Replicas = int(deployment.Spec.Replicas)
	if Replicas == 0 {
		fmt.Println("[yaml ERROR]Replicas==0")
		return nil, nil
	}
	//根据template获得template hash
	templateHash := strconv.Itoa(int(HASH.HASH([]byte(deployment.Metadata.Name + strconv.Itoa(int(deployment.Spec.Replicas))))))

	Pods := make(map[string]*entity.Pod, Replicas)
	for i := 0; i < Replicas; i++ {
		//创建replicas份Pod
		pod := &entity.Pod{}
		pod.Kind = "Pod"
		pod.Metadata = deployment.Spec.Template.Metadata
		uid := UUID.UUID()
		pod.Metadata.Uid = uid
		//组合产生Deployment pod的名字
		pod.Metadata.Name = deployment.Metadata.Name + "-" + templateHash + "-" + uid[:5]

		pod.Spec = deployment.Spec.Template.Spec
		Pods[pod.Metadata.Name+pod.Metadata.Uid] = pod

		////使用模板启动时，之后的replica port使用递增号,更新container中的端口映射信息
		//for j, con := range pod.Spec.Containers {
		//	for m, port := range con.Ports {
		//		//port.ContainerPort = oldportToNewport(port.ContainerPort, i)
		//		if port.HostPort != "" {
		//			port.HostPort = oldportToNewport(port.HostPort, i)
		//		} else {
		//			port.HostPort = strconv.Itoa(PORT.GetFreePort())
		//		}
		//		pod.Spec.Containers[j].Ports[m] = port
		//	}
		//}

	}
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
	//deployment信息写入etcd
	deployment.Status.StartTime = time.Now()
	DeploymentData, err := json.Marshal(deployment)
	if err != nil {
		fmt.Println("convert DeploymentYaml error!")
	}
	err = etcdctl.Put(cli, "Deployment/"+deployment.Metadata.Name, string(DeploymentData))
	if err != nil {
		fmt.Println("write DeploymentYaml to etcd error!")
		return nil, err
	}
	//pod信息写入etcd
	podsList := []*entity.Pod{}
	for _, pod := range Pods {
		podsList = append(podsList, pod)
		PodData, err := json.Marshal(pod)
		err = etcdctl.Put(cli, "Pod/"+pod.Metadata.Name, string(PodData))
		if err != nil {
			fmt.Println("write PodData to etcd error!")
			return nil, err
		}
	}
	//写入成功
	fmt.Println("*************write deployment success!************")
	return podsList, err
}

// DeleteDeployment kubectl调用后删除etcd中的deployment信息
func DeleteDeployment(DeploymentName string) error {
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
	//删除deployment
	//解析deployment的内容
	deploymentDetail := entity.Deployment{}
	deployment := make([]byte, 0)
	deployments, _ := etcdctl.Get(cli, "Deployment/"+DeploymentName)
	if len(deployments.Kvs) != 0 {
		deployment = deployments.Kvs[0].Value
		err = json.Unmarshal(deployment, deploymentDetail)
	}
	etcdctl.Delete(cli, "Deployment/"+DeploymentName)
	//查询etcd中和deployment相关的pod
	PodsData, err := cli.Get(context.Background(), "Pod/"+deploymentDetail.Metadata.Name, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("get from etcd failed, err:", err)
		return err
	}
	//通知Node删除和deployment相关的pod，更新etcd
	Pods := make(map[string]entity.Pod)
	for _, value := range PodsData.Kvs {
		pod := entity.Pod{}
		err = json.Unmarshal(value.Value, &pod)
		Pods[pod.Metadata.Name] = pod
	}
	//TODO 通知Node删除Pod,利用hostip

	log.Print("Node删除成功后删除etcd信息")
	//Node删除成功后删除etcd信息
	for _, pod := range Pods {
		podpath := "Pod/" + pod.Metadata.Name
		err := etcdctl.Delete(cli, podpath)
		if err != nil {
			log.Print("delete " + podpath + "failed!")
			return err
		}
	}
	return err
}

// GetPodsBydeployment Pod命名：deploymentname(nginx-deployment)+templet对应HASH(9594276)+PodUID前5位
func GetPodsBydeployment(deployment string) []entity.Pod {
	//TODO 默认了能在etcd中查询到的pod都是可用状态，属于replica，是否需要更新？
	PodsData, _ := etcdctl.EtcdGetWithPrefix("Pod/" + deployment)
	var Pods []entity.Pod
	for _, PodData := range PodsData.Kvs {
		var pod entity.Pod
		err := json.Unmarshal(PodData.Value, &pod)
		if err != nil {
			log.PrintE("GetPodsBydeployment Unmarshal Pod error")
			return nil
		}
		Pods = append(Pods, pod)
	}
	return Pods
}

// MonitorDeployment 另开线程中运行，持续检查deployment运行状态，并进行扩缩
func MonitorDeployment() error {
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
	//获取当前所有的deployment
	Deployments := make(map[string]entity.Deployment)
	deploymentByte := make([]byte, 0)
	deployments, _ := cli.Get(context.Background(), "Deployment/", clientv3.WithPrefix())
	for _, value := range deployments.Kvs {
		deploymentDetail := entity.Deployment{}
		deploymentByte = value.Value
		err = json.Unmarshal(deploymentByte, &deploymentDetail)
		Deployments[deploymentDetail.Metadata.Name] = deploymentDetail
	}

	//获取其中的pod状态信息，failed则通知重新创建
	//不满足replica要求则进行增删
	FailedPods := make(map[string]entity.Pod)
	for _, deployment := range Deployments {
		PodsData, err := cli.Get(context.Background(), "Pod/"+deployment.Metadata.Name, clientv3.WithPrefix())
		if err != nil {
			fmt.Println("get from etcd failed, err:", err)
			return err
		}
		moreNum, fewerNum := 0, 0
		//大于replica,删除多余的pod
		if int32(len(PodsData.Kvs)) > deployment.Spec.Replicas {
			fewerNum = len(PodsData.Kvs) - int(deployment.Spec.Replicas)
		}
		//小于replica,补充不足的的pod
		if int32(len(PodsData.Kvs)) < deployment.Spec.Replicas {
			moreNum = int(deployment.Spec.Replicas) - len(PodsData.Kvs)
		}
		for _, value := range PodsData.Kvs {
			pod := entity.Pod{}
			err = json.Unmarshal(value.Value, &pod)
			if fewerNum != 0 {
				//删除failed pod
				_, err := cli.Delete(context.Background(), "Pod/"+pod.Metadata.Name)
				if err != nil {
					fmt.Println("[MonitorDeployment]delete failed pod fail!")
					return err
				}
				fewerNum--
			} else {
				//如果Pod状态为Failed且fewerNum=0则记录
				if pod.Status.Phase == entity.Failed {
					FailedPods[pod.Metadata.Name] = pod
				}
			}
		}
		//deployment.replica有更新，需要增加moreNum个 pod
		//TODO 新增etcd的pod信息后，如何通知Node的kubelet更新？
		if moreNum != 0 {
			for i := 0; i <= moreNum; i++ {
				//创建replicas份Pod
				pod := &entity.Pod{}
				pod.Metadata = deployment.Spec.Template.Metadata
				pod.Metadata.Uid = UUID.UUID()
				//根据template获得template hash
				templateHash := strconv.Itoa(int(HASH.HASH([]byte(deployment.Metadata.Name + strconv.Itoa(int(deployment.Spec.Replicas))))))
				//组合产生Deployment pod的名字
				pod.Metadata.Name = deployment.Metadata.Name + "-" + templateHash + "-" + pod.Metadata.Uid[:5]
				pod.Spec = deployment.Spec.Template.Spec
				PodData, err := json.Marshal(pod)
				err = etcdctl.Put(cli, "Pod/"+pod.Metadata.Name, string(PodData))
				if err != nil {
					fmt.Println("[MonitorDeployment]write more PodData to etcd error!")
					return err
				}
			}
		}
	}
	cli.Close()
	//TODO 通知Node补充更新所有的FailedPod

	return nil
}
