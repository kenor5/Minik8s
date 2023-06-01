package ControllerManager

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
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
	//templateHash := strconv.Itoa(int(HASH.HASH([]byte(deployment.Metadata.Name + strconv.Itoa(int(deployment.Spec.Replicas))))))
	templateHash := strconv.Itoa(int(HASH.HASH([]byte(deployment.Metadata.Name))))
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

		pod.Status.Phase = entity.Pending
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
			log.PrintE("close etcdClient error!")
		}
	}(cli)
	//deployment信息写入etcd
	deployment.Status.StartTime = time.Now()
	DeploymentData, err := json.Marshal(deployment)
	if err != nil {
		fmt.Println("convert DeploymentYaml error!")
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
	err = etcdctl.Put(cli, "Deployment/"+deployment.Metadata.Name, string(DeploymentData))
	if err != nil {
		fmt.Println("write DeploymentYaml to etcd error!")
		return nil, err
	}
	//写入成功
	fmt.Println("*************write deployment success!************")
	return podsList, err
}

// DeleteDeployment kubectl调用后删除etcd中的deployment信息
func DeleteDeployment(DeploymentName string) error {
	log.Printf("[DeleteDeployment]Beging delete Deployment:%s", DeploymentName)
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connect error")
	}
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			log.PrintE("close etcdClient error!")
		}
	}(cli)
	//删除deployment
	//解析deployment的内容
	deploymentDetail := &entity.Deployment{}
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
		log.PrintE("get from etcd failed, err:", err)
		return err
	}
	//通知Node删除和deployment相关的pod，更新etcd
	Pods := make(map[string]entity.Pod)
	for _, value := range PodsData.Kvs {
		pod := entity.Pod{}
		err = json.Unmarshal(value.Value, &pod)
		Pods[pod.Metadata.Name] = pod
	}

	log.Print("Node删除成功后删除etcd信息")
	//Node删除成功后删除etcd信息,借助Pod更新机制
	for _, pod := range Pods {
		podpath := "Pod/" + pod.Metadata.Name
		pod.Status.Phase = entity.Succeed
		podData, _ := json.Marshal(pod)
		err := etcdctl.Put(cli, "Pod/"+pod.Metadata.Name, string(podData))
		if err != nil {
			log.PrintE("delete " + podpath + " failed!")
			return err
		}
	}
	//for _, pod := range Pods {
	//	podpath := "Pod/" + pod.Metadata.Name
	//	err := etcdctl.Delete(cli, podpath)
	//	if err != nil {
	//		log.PrintE("delete " + podpath + "failed!")
	//		return err
	//	}
	//}
	return err
}

// GetPodsBydeployment Pod命名：deploymentname(nginx-deployment)+templet对应HASH(9594276)+PodUID前5位
func GetPodsBydeployment(deployment string) []entity.Pod {
	PodsData, _ := etcdctl.EtcdGetWithPrefix("Pod/" + deployment)
	var Pods []entity.Pod
	for _, PodData := range PodsData.Kvs {
		var pod entity.Pod
		err := json.Unmarshal(PodData.Value, &pod)
		if err != nil {
			log.PrintE("GetPodsBydeployment Unmarshal Pod error")
			return nil
		}
		if pod.Status.Phase == entity.Running {
			Pods = append(Pods, pod)
		}
	}
	return Pods
}

func BeginMonitorDeployment() {
	//ticker := time.NewTicker(configs.MonitorDeploymentTime * time.Second) //每10s使用
	//for range ticker.C {
	//	err := MonitorDeployment()
	//	if err != nil {
	//		log.PrintE(err)
	//		continue
	//	}
	//}
	for true {
		err := WatchMonitorDeployment()
		if err != nil {
			log.PrintE(err)
			continue
		}
	}
}
func WatchMonitorDeployment() error {
	log.Printf("[ApiServer]Begin Monitor Deployment")
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connect error")
	}
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			log.PrintE("close etcdClient error!")
		}
	}(cli)
	watchCh := cli.Watch(context.Background(), "Deployment/", clientv3.WithPrefix())
	// 接收 key 的变化事件
	for resp := range watchCh {
		for _, event := range resp.Events {
			if event.Type == mvccpb.PUT { // 如果 key 值被更新了
				fmt.Printf("key=%s updatedvalue=%s\n", event.Kv.Key, event.Kv.Value)
				deployment := &entity.Deployment{}
				json.Unmarshal(event.Kv.Value, deployment)
				go func() {
					err := CheckDeploymentPod(*deployment)
					if err != nil {
						log.PrintfE("Check deployment %s error", deployment.Metadata.Name)
					}
				}()
			} else if event.Type == mvccpb.DELETE { // 如果 key 值被删除了
				fmt.Printf("key=%s deleted=true\n", event.Kv.Key)
			}
		}
	}

	return err
}

func CheckDeploymentPod(deployment entity.Deployment) error {
	PodsData, err := etcdctl.EtcdGetWithPrefix("Pod/" + deployment.Metadata.Name)
	if err != nil {
		fmt.Println("get from etcd failed, err:", err)
		return err
	}
	moreNum, fewerNum := 0, 0
	nowReplica := 0
	for _, value := range PodsData.Kvs {
		pod := entity.Pod{}
		err = json.Unmarshal(value.Value, &pod)
		if pod.Status.Phase == entity.Running || pod.Status.Phase == entity.Pending {
			nowReplica++
		}
	}
	log.Printf("[MonitorDeployment]now replica:%d need replica:%d", nowReplica, deployment.Spec.Replicas)
	deployment.Status.Replicas = int32(nowReplica)
	//大于replica,删除多余的pod
	if deployment.Status.Replicas > deployment.Spec.Replicas {
		//fewerNum = len(PodsData.Kvs) - int(deployment.Spec.Replicas)
		fewerNum = int(deployment.Status.Replicas - deployment.Spec.Replicas)
		log.Printf("[MonitorDeployment]fewerNum=%d", fewerNum)
	}
	//小于replica,补充不足的的pod
	if deployment.Status.Replicas < deployment.Spec.Replicas {
		//moreNum = int(deployment.Spec.Replicas) - len(PodsData.Kvs)
		moreNum = int(deployment.Spec.Replicas - deployment.Status.Replicas)
		log.Printf("[MonitorDeployment]moreNum=%d", moreNum)
	}
	for _, value := range PodsData.Kvs {
		pod := entity.Pod{}
		err = json.Unmarshal(value.Value, &pod)
		if pod.Status.Phase == entity.Succeed {
			continue
		}
		if fewerNum > 0 {
			//TODO：删除failed pod，改为succeed，借助Pod更新机制
			log.Printf("[MonitorDeployment]%s need sub %d replica", deployment.Metadata.Name, fewerNum)
			pod.Status.Phase = entity.Succeed
			PodsData, _ := json.Marshal(pod)
			err := etcdctl.EtcdPut("Pod/"+pod.Metadata.Name, string(PodsData))
			if err != nil {
				fmt.Println("[MonitorDeployment]delete failed pod fail!")
				return err
			}
			fewerNum--
		} else {
			//如果Pod状态为Failed且fewerNum=0则记录
			break
		}
	}
	//deployment.replica有更新，需要增加moreNum个 pod
	//TODO 新增etcd的pod信息后，如何通知Node的kubelet更新？借助Pod更新机制
	if moreNum > 0 {
		log.Printf("[MonitorDeployment]%s need add %d replica", deployment.Metadata.Name, moreNum)
		for i := 0; i <= moreNum; i++ {
			//创建replicas份Pod
			pod := &entity.Pod{}
			pod.Metadata = deployment.Spec.Template.Metadata
			pod.Metadata.Uid = UUID.UUID()
			//根据template获得template hash
			//templateHash := strconv.Itoa(int(HASH.HASH([]byte(deployment.Metadata.Name + strconv.Itoa(int(deployment.Spec.Replicas))))))
			templateHash := strconv.Itoa(int(HASH.HASH([]byte(deployment.Metadata.Name))))

			//组合产生Deployment pod的名字
			pod.Metadata.Name = deployment.Metadata.Name + "-" + templateHash + "-" + pod.Metadata.Uid[:5]
			pod.Spec = deployment.Spec.Template.Spec
			pod.Status.Phase = entity.Pending
			PodData, err := json.Marshal(pod)
			err = etcdctl.EtcdPut("Pod/"+pod.Metadata.Name, string(PodData))
			if err != nil {
				fmt.Println("[MonitorDepployment]write more PodData to etcd error!")
				return err
			}
		}
	}
	return err
}

// MonitorDeployment 另开线程中运行，持续轮询检查deployment运行状态，并进行扩缩
func MonitorDeployment() error {
	log.Printf("[ApiServer]Begin Monitor Deployment")
	//获取当前所有的deployment
	Deployments := make(map[string]entity.Deployment)
	deploymentByte := make([]byte, 0)
	deployments, _ := etcdctl.EtcdGetWithPrefix("Deployment/")
	for _, value := range deployments.Kvs {
		deploymentDetail := entity.Deployment{}
		deploymentByte = value.Value
		json.Unmarshal(deploymentByte, &deploymentDetail)
		Deployments[deploymentDetail.Metadata.Name] = deploymentDetail
	}

	//获取其中的pod状态信息，failed则通知重新创建
	//不满足replica要求则进行增删
	for _, deployment := range Deployments {
		PodsData, err := etcdctl.EtcdGetWithPrefix("Pod/" + deployment.Metadata.Name)
		if err != nil {
			fmt.Println("get from etcd failed, err:", err)
			return err
		}
		moreNum, fewerNum := 0, 0
		nowReplica := 0
		for _, value := range PodsData.Kvs {
			pod := entity.Pod{}
			err = json.Unmarshal(value.Value, &pod)
			if pod.Status.Phase == entity.Running || pod.Status.Phase == entity.Pending {
				nowReplica++
			}
		}
		log.Printf("[MonitorDeployment]now replica:%d need replica:%d", nowReplica, deployment.Spec.Replicas)
		deployment.Status.Replicas = int32(nowReplica)
		//大于replica,删除多余的pod
		if deployment.Status.Replicas > deployment.Spec.Replicas {
			//fewerNum = len(PodsData.Kvs) - int(deployment.Spec.Replicas)
			fewerNum = int(deployment.Status.Replicas - deployment.Spec.Replicas)
			log.Printf("[MonitorDeployment]fewerNum=%d", fewerNum)
		}
		//小于replica,补充不足的的pod
		if deployment.Status.Replicas < deployment.Spec.Replicas {
			//moreNum = int(deployment.Spec.Replicas) - len(PodsData.Kvs)
			moreNum = int(deployment.Spec.Replicas - deployment.Status.Replicas)
			log.Printf("[MonitorDeployment]moreNum=%d", moreNum)
		}
		for _, value := range PodsData.Kvs {
			pod := entity.Pod{}
			err = json.Unmarshal(value.Value, &pod)
			if pod.Status.Phase == entity.Succeed {
				continue
			}
			if fewerNum > 0 {
				//TODO：删除failed pod，改为succeed，借助Pod更新机制
				log.Printf("[MonitorDeployment]%s need sub %d replica", deployment.Metadata.Name, fewerNum)
				pod.Status.Phase = entity.Succeed
				PodsData, _ := json.Marshal(pod)
				err := etcdctl.EtcdPut("Pod/"+pod.Metadata.Name, string(PodsData))
				if err != nil {
					fmt.Println("[MonitorDeployment]delete failed pod fail!")
					return err
				}
				fewerNum--
			} else {
				//如果Pod状态为Failed且fewerNum=0则记录
				break
			}
		}
		//deployment.replica有更新，需要增加moreNum个 pod
		//TODO 新增etcd的pod信息后，如何通知Node的kubelet更新？借助Pod更新机制
		if moreNum > 0 {
			log.Printf("[MonitorDeployment]%s need add %d replica", deployment.Metadata.Name, moreNum)
			for i := 0; i <= moreNum; i++ {
				//创建replicas份Pod
				pod := &entity.Pod{}
				pod.Metadata = deployment.Spec.Template.Metadata
				pod.Metadata.Uid = UUID.UUID()
				//根据template获得template hash
				//templateHash := strconv.Itoa(int(HASH.HASH([]byte(deployment.Metadata.Name + strconv.Itoa(int(deployment.Spec.Replicas))))))
				templateHash := strconv.Itoa(int(HASH.HASH([]byte(deployment.Metadata.Name))))

				//组合产生Deployment pod的名字
				pod.Metadata.Name = deployment.Metadata.Name + "-" + templateHash + "-" + pod.Metadata.Uid[:5]
				pod.Spec = deployment.Spec.Template.Spec
				pod.Status.Phase = entity.Pending
				PodData, err := json.Marshal(pod)
				err = etcdctl.EtcdPut("Pod/"+pod.Metadata.Name, string(PodData))
				if err != nil {
					fmt.Println("[MonitorDepployment]write more PodData to etcd error!")
					return err
				}
			}
		}
	}
	//TODO:通知Node补充更新所有的FailedPod,借助Pod更新机制

	return nil
}
