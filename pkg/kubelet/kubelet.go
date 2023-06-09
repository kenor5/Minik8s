package kubelet

import (
	// "context"
	// "encoding/json"
	// "fmt"

	"context"
	"encoding/json"
	"fmt"
	"minik8s/tools/log"
	"minik8s/tools/yamlParser"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"

	"minik8s/configs"
	"minik8s/entity"
	"minik8s/pkg/kubelet/pod/PodManager"
	"minik8s/pkg/kubelet/pod/podfunc"
	"os"

	kp "minik8s/pkg/kube_proxy"
	pb "minik8s/pkg/proto"
	"regexp"
	"strconv"
	"sync"
	"time"

	// "net"

	"minik8s/pkg/kubelet/client"
	"minik8s/pkg/kubelet/container/ContainerManager"
	"minik8s/tools/network"

	//"minik8s/pkg/kubelet/pod/PodManager"

	// "github.com/docker/docker/libnetwork/drivers/host"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/**************************************************************************
************************    Kubelet主结构    *******************************
***************************************************************************/
type Kubelet struct {
	hostName string
	hostIp   string

	lock             sync.Locker
	connToApiServer  pb.ApiServerKubeletServiceClient // kubelet连接到apiserver的conn
	podManger        PodManager.PodManager            //存储在内存中的pod信息
	containerManager *ContainerManager.ContainerManager
}

var kubelet *Kubelet
var kubeProxy *kp.KubeProxy

// newKubelet creates a new Kubelet object.
func newKubelet() *Kubelet {
	newKubelet := &Kubelet{
		lock:             &sync.RWMutex{},
		podManger:        PodManager.NewPodManager(),
		containerManager: ContainerManager.NewContainerManager(),
	}
	apiserver_url := configs.ApiServerUrl + configs.GrpcPort
	newKubelet.connToApiServer, _ = ConnectToApiServer(apiserver_url)
	// 获取主机名和主机IP
	hostname, _ := os.Hostname()
	newKubelet.hostName = hostname
	IP, err := network.GetNetInterfaceIPv4Addr(configs.NetInterface)
	if err != nil {
		log.PrintE("fail to get hostIP!")
	}
	newKubelet.hostIp = IP
	return newKubelet
}

func KubeletObject() *Kubelet {
	if kubelet == nil {
		kubelet = newKubelet()
	}
	return kubelet
}

func KubeProxyObject() *kp.KubeProxy {
	var err error
	if kubeProxy == nil {
		kubeProxy, err = kp.NewKubeProxy()
		if err != nil {
			log.PrintE("error when creating kubeproxy")
			return nil
		}
		return kubeProxy
	}
	return kubeProxy
}

func (kl *Kubelet) CreatePod(pod *entity.Pod) error {
	//kl.lock.Lock()
	//defer kl.lock.Unlock()
	// 实际创建Pod,IP等信息在这里更新进Pod.Status中
	pod.Status.HostIp = kl.hostIp
	pod.Spec.NodeName = kl.hostName
	ContainerIds, err := podfunc.CreatePod(pod)
	if err != nil {
		log.PrintE(err)
		return err
	}

	log.PrintS("1")
	// 维护ContainerRuntimeManager
	kl.podManger.AddPod(*pod)

	log.PrintS("2")

	for _, ContainerId := range ContainerIds {
		log.PrintS(ContainerId)
		kubelet.podManger.AddContainerToPod(ContainerId, *pod)
	}
	kl.containerManager.SetContainerIDsByPodName(pod, ContainerIds)

	log.PrintS("3")

	// 更新PodStatus
	log.PrintS("[Kubelet] CreatePod success,Begin update Pod")
	client.UpdatePodStatus(kubelet.connToApiServer, pod)

	log.PrintS("4")
	return nil
}

func (kl *Kubelet) DeletePod(pod *entity.Pod) error {
	// 获取Pod中所有的ContainerId并且删除该映射
	containerIds := kubelet.containerManager.GetContainerIDsByPodName(pod.Metadata.Name)
	kl.containerManager.DeletePodNameToContainerIds(pod.Metadata.Name)

	log.Print("containerIds: %s\n", containerIds)
	// 实际停止并删除Pod中的所有容器
	podfunc.DeletePod(containerIds)
	//kl.podManger.DeletePod(pod)
	// 更新Pod的状态
	pod.Status.Phase = entity.Succeed
	kl.podManger.DeletePod(*pod)
	log.Print("[Kubelet] DeletePod success,Begin update Pod")
	client.UpdatePodStatus(kubelet.connToApiServer, pod)
	return nil
}

func (kl *Kubelet) AddPod(pod *entity.Pod) error {
	//更新元数据
	//kl.podManger.AddPod(pod)
	pod.Status.Phase = entity.Running

	//启动沙箱容器和pod.spec.containers中的容器
	if _, err := podfunc.CreatePod(pod); err != nil {
		// pod.Status.Phase = entity.Failed
		return err
	}

	return nil
}

// func (kl *Kubelet) GetPodByName(namespace string, name string) (*entity.Pod, bool) {
// 	pm, ok := kl.podManger.GetPodByName(namespace, name)
// 	return pm, ok
// }

func (kl *Kubelet) RegisterNode() error {
	// registerNodeRequest := &pb.RegisterNodeRequest{
	// 	NodeName:   kl.hostName,
	// 	KubeletUrl: kl.hostIp + configs.KubeletGrpcPort,
	// }
	newNode := &entity.Node{}
	yamlParser.ParseYaml(newNode, configs.NodeConfig)
	kl.hostName = newNode.Name
	Podsbyte, err := client.RegisterNode(kl.connToApiServer, kl.hostName, kl.hostIp)
	//注册返回Pod信息并初始化本地信息,用APIServer中接口，初始化时获得HostIP 为自己的pod列表
	for _, podByte := range Podsbyte {
		pod := entity.Pod{}
		err := json.Unmarshal(podByte, &pod)
		kl.podManger.AddPod(pod)
		err = kl.getContainersByPodname(pod)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (kl *Kubelet) getContainersByPodname(pod entity.Pod) error {
	podname := pod.Metadata.Name
	ctx := context.Background()
	cli, err := dockerclient.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// 获取运行中的容器列表
	containers, err := cli.ContainerList(ctx, dockertypes.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	// 创建一个 filter 查询条件
	filter := filters.NewArgs()
	filter.Add("name", podname)

	// 获取容器列表
	containers, err = cli.ContainerList(context.Background(), dockertypes.ContainerListOptions{
		All:     false,
		Filters: filter,
	})
	if err != nil {
		panic(err)
	}

	// 该map返回Pod中的ContainerID
	ContainerIDMap := []string{}

	// 遍历容器列表输出 ID
	for _, container := range containers {
		kl.podManger.AddContainerToPod(container.ID, pod)
		ContainerIDMap = append(ContainerIDMap, container.ID)
	}

	kl.containerManager.SetContainerIDsByPodName(&pod, ContainerIDMap)

	return err

}

func ConnectToApiServer(apiserver_url string) (pb.ApiServerKubeletServiceClient, error) {
	dial, err := grpc.Dial(apiserver_url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.PrintE(err)
		return nil, err
	}
	// defer dial.Close()

	conn := pb.NewApiServerKubeletServiceClient(dial)
	return conn, err
}

// 在一个单独的线程中运行，监控pod状态
func (kl *Kubelet) monitorPods() {
	//kl.lock.Lock()
	//defer kl.lock.Unlock()
	log.Printf("[Kubelet]Begin Monitor Pod")
	var FailedPods []*entity.Pod
	var SucceedPods []*entity.Pod
	cli, _ := dockerclient.NewClientWithOpts(
		dockerclient.FromEnv,
		dockerclient.WithAPIVersionNegotiation(),
	)
	defer cli.Close()
	//通知API server更新Pod状态
	//遍历当前Node的pod
	log.Print("1")
	Pods := kl.podManger.GetPods()
	for _, pod := range Pods {
		if pod.Status.Phase == entity.Failed {
			FailedPods = append(FailedPods, &pod)
			//client.UpdatePodStatus(kl.connToApiServer, pod)
			continue
		}
		//Succeed为主动关闭的pod
		if pod.Status.Phase == entity.Succeed {
			SucceedPods = append(SucceedPods, &pod)
			//client.UpdatePodStatus(kl.connToApiServer, pod)
			continue
		}
		exitCodeReg, _ := regexp.Compile(`\(\d+\)`)
		isFailed := false
		isFinished := false
		//通过ContainerID检查Running Pod的container运行状态
		containers := kl.podManger.GetContainersByPod(pod)

		log.Print("len(containers): ", len(containers))

		if len(containers) > 0 {
			log.Print("2")
			for _, containerId := range containers {
				if isFinished || isFailed {
					//当有Container 退出或者终止，则记录pod为Failed
					break
				}
				containers, err := cli.ContainerList(context.Background(), dockertypes.ContainerListOptions{
					All: true,
					Filters: filters.NewArgs(
						filters.Arg("id", containerId),
					),
				})
				if err != nil {
					fmt.Printf("fail to query container %v's status: %v", containerId, err)
					continue
				}

				log.Print("3")

				container := containers[0]
				switch container.State {
				case "exited":
					m := exitCodeReg.FindString(container.Status)
					exitCode, err := strconv.Atoi(m[1 : len(m)-1])
					if err != nil {
						fmt.Printf("fail to parse container %v's exit code: %v", containerId, err)
					}
					if exitCode != 0 {
						isFailed = true
						FailedPods = append(FailedPods, &pod)
					}
				case "dead":
					isFailed = true
					FailedPods = append(FailedPods, &pod)
				default:
					isFinished = false
				}
			}
		} else {
			pod.Status.Phase = entity.Failed
		}

	}

    log.Print("4")
	
	//删除Succeed Pod的pod和container
	for _, pod := range SucceedPods {
		conatainers := kl.podManger.GetContainersByPod(*pod)
		err := podfunc.DeletePod(conatainers)
		if err != nil {
			fmt.Printf("delete containerId %v error", conatainers[0])
			return
		}
		kl.podManger.DeleteContainersByPod(*pod)
		kl.podManger.DeletePod(*pod)
	}

	//删除FailedPod的其余container，并重新创建Pod
	for _, pod := range FailedPods {
		err := client.UpdatePodStatus(kubelet.connToApiServer, pod)
		if err != nil {
			log.PrintfE("[monitorPods]Update Pod %s status erroe", pod.Metadata.Name)
			continue
		}
		conatainers := kl.podManger.GetContainersByPod(*pod)
		err = podfunc.DeletePod(conatainers)
		if err != nil {
			fmt.Printf("delete container of Pod %v error", pod.Metadata.Name)
			return
		}
		kl.podManger.DeleteContainersByPod(*pod)
	}
	//kl.lock.Unlock()
	//重新创建FailedPod
	for _, pod := range FailedPods {
		err := kl.CreatePod(pod)
		if err != nil {
			fmt.Printf("create Pod %v error!", pod.Metadata.Name)
			return
		}

	}

}

// 每30s检查一次本地运行Pod状态
// 使用 go beginMonitor()开始执行
func (kl *Kubelet) BeginMonitor() {

	ticker := time.NewTicker(30 * time.Second) //每30s使用
	for range ticker.C {
		kl.monitorPods()
	}
	/*go func() {
	for range time.Tick(time.Second * monitorInterval) {
	kubelet.monitorPods()
	}
	}()*/
}

func (kl *Kubelet) ApplyDns(dns *entity.Dns) error {
	err := CreateDns(dns)
	if err != nil {
		return err
	}
	return nil
}

func (kl *Kubelet) DeleteDns(dnsName string) error {
	//TODO
	return nil
}
