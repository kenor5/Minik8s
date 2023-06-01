package kube_proxy

import (
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/kube_proxy/iptable"
	"minik8s/tools/log"
	"net"
	"strings"
)

type KubeProxy struct {
	IptableClient  *iptable.IpTable
	ServiceManager *ServiceManager
}

func NewKubeProxy() (*KubeProxy, error) {

	fmt.Println("creating new proxy")
	iptableCli, err := iptable.NewClient()
	if err != nil {
		log.PrintE("error when create iptable client")
		log.PrintE(err)
		log.PrintW("gjl notes: if '/run/xtables.lock: Permission denied', pleast run in root")
		return nil, nil
	}

	err = iptableCli.Init()
	if err != nil {
		log.PrintE("error when create init iptable client")
		log.PrintE(err)
		log.PrintW("gjl notes: if '/run/xtables.lock: Permission denied', pleast run in root")
		return nil, err
	}
	log.PrintS("success in create kp")

	return &KubeProxy{
		IptableClient:  iptableCli,
		ServiceManager: NewSvcManager(),
	}, nil
}

func (kp *KubeProxy) NewService(service *entity.Service, podNames []string, podIps []string) error {
	servicePorts := service.Spec.Ports
	serviceName := service.Metadata.Name
	podLen := len(podNames)
	// 需要再内存和 iptable 中分别加上相应的信息

	// 对于每一个端口，需要创建一个 service chain， 并为这个 chain 添加相应的跳转规则
	for _, port := range servicePorts {
		svcChainName := kp.IptableClient.CreateServiceChain()
		// [内存] 加入 svcName -> svcChain 的映射
		kp.ServiceManager.AddServiceChain(serviceName, svcChainName, &port)
		kp.ServiceManager.AddClusterIp(serviceName, service.Spec.ClusterIP)

		for i := 0; i < podLen; i++ {
			//podName := podNames[i]
			podIp := podIps[i]
			podChainName := kp.IptableClient.CreatePodChain()

			// [iptable] pod chain 跳转到具体pod
			// TODO check if pod belongs to this Node
			err := kp.IptableClient.AddPodRules(podChainName, podIp, uint32(port.TargetPort))
			if err != nil {
				return err
			}

			probablity := 1.0/float64(podLen-i)
			// [内存] service chain 跳转到具体 pod chain
			// err = kp.IptableClient.ApplyPodRules(svcChainName, podChainName, podLen, i)
			err = kp.IptableClient.ApplyPodRules2(svcChainName, podChainName, probablity)
			if err != nil {
				return err
			}

			kp.ServiceManager.AddPodChain(svcChainName, podChainName, podNames[i], probablity)

		}

		// [iptable] KUBE-SERVICE 跳转到具体的service
		err := kp.IptableClient.AddServiceRules(service.Spec.ClusterIP, svcChainName, uint32(port.Port))
		if err != nil {
			return err
		}
	}


	return nil
}

func (kp *KubeProxy) RemoveService(serviceName string) error {
	// 需要删除 service manager 中的数据和 iptable 中的路由数据
	serviceChains := kp.ServiceManager.GetServiceChains(serviceName)
	clusterIp := kp.ServiceManager.GetClusterIp(serviceName)

	for _, svc := range serviceChains {
		// 先删除 iptable 里的路由规则
		err := kp.IptableClient.RemoveServiceChain(svc.ChainName, clusterIp, uint32(svc.Ports.Port))
		if err != nil {
			log.PrintE(err)
			return err
		}

		podChains := kp.ServiceManager.GetPodChains(svc.ChainName)
		log.PrintE("pod chain %v", podChains)
		for _, podChain := range podChains {
			err := kp.IptableClient.RemovePodChain(podChain.ChainName)
			if err != nil {
				log.PrintE(err)
				return err
			}
		}
		// 在删除内存里的 podchain 信息
		kp.ServiceManager.RemoveServiceChain(svc.ChainName)
	}

	kp.ServiceManager.RemoveServiceChain(serviceName)

	return nil
}

// targetPort 就是service yaml 文件中的targetPort
func (kp *KubeProxy) AddPod2Service(svcName string, podName string, podIp string, targetPort uint32) error  {
	log.Print("begin add pod 2 service")

	// [iptable] 在KUBE-SVC-XXX链中添加规则
	if !kp.ServiceManager.ExistServiceChain(svcName) {
		log.PrintE("no such service chain")
	}

	svcChains := kp.ServiceManager.GetServiceChains(svcName)
	// 把原来的记录都删掉
	for _, svcChain := range svcChains {
		podChainNames := kp.ServiceManager.GetPodChains(svcChain.ChainName)
		for _,podChainName := range podChainNames {
			err := kp.IptableClient.RemovePodRules2(svcChain.ChainName, podChainName.ChainName, podChainName.probability)
			if err != nil {
				log.PrintE(err)
				return err
			}
		}
		log.Print("original pod len: ", len(podChainNames))
	}
	// 创建KUBE-SEP-XXX链
	podChainName := kp.IptableClient.CreatePodChain()

	// [内存] 在servicechainname2podchain 中添加一条记录
	for _, svcChain := range svcChains {

		// 在iptable中写回原来的记录和新的记录
		kp.ServiceManager.AddPodChain(svcChain.ChainName, podChainName, podName, 0.5)

		err := kp.IptableClient.AddPodRules(podChainName, podIp, targetPort)
		if err != nil {
			log.PrintE(err)
			return err
		}

		podChains := kp.ServiceManager.GetPodChains(svcChain.ChainName)
		podLen := len(podChains)
		for i, podChain := range podChains {
			probability := 1.0/float64(podLen-i)
			kp.ServiceManager.ServiceChainName2PodChain[svcChain.ChainName][i].probability = probability
			kp.IptableClient.ApplyPodRules2(svcChain.ChainName, podChain.ChainName, probability)
		}
		log.Print("cur pod len: ", podLen)
	}


	return nil
}

func (kp *KubeProxy) RemovePodFromService(svcName string, podName string) error {
	// 从 KUBE-SVC-XXX 链中删除引用

	// 删除 KUBE-SEP-xxx 链

	// [memory] 从servicechainname2podchain中删除一条记录

	log.Print("begin remove pod from service")

	// [iptable] 在KUBE-SVC-XXX链中添加规则
	if !kp.ServiceManager.ExistServiceChain(svcName) {
		log.PrintE("no such service chain")
	}
	svcChains := kp.ServiceManager.GetServiceChains(svcName)
	// 把原来的记录都删掉
	for _, svcChain := range svcChains {
		podChainNames := kp.ServiceManager.GetPodChains(svcChain.ChainName)
		for _,podChainName := range podChainNames {
			err := kp.IptableClient.RemovePodRules2(svcChain.ChainName, podChainName.ChainName, podChainName.probability)
			if err != nil {
				log.PrintE(err)
			}
		}
		
		log.Print("original pod len: %d", len(podChainNames))
	}

	
	var podChainName string
	for _, svcChain := range svcChains {
		podChainName = kp.ServiceManager.GetPodChainNameByPodName(svcChain.ChainName, podName)
		if podChainName != "" {
			break
		}
	}

	if podChainName == "" {
		log.PrintE("no such pod chain")
		return nil
	}
	// [内存] 在servicechainname2podchain 中添加一条记录
	for _, svcChain := range svcChains {
		
		// 在iptable中写回原来的记录和新的记录
		kp.ServiceManager.RemovePodChain(svcChain.ChainName, podChainName)

		podChains := kp.ServiceManager.GetPodChains(svcChain.ChainName)
		podLen := len(podChains)
		for i, podChain := range podChains {
			probability := 1.0/float64(podLen-i)
			kp.ServiceManager.ServiceChainName2PodChain[svcChain.ChainName][i].probability = probability
			kp.IptableClient.ApplyPodRules2(svcChain.ChainName, podChain.ChainName, probability)
		}
		log.Print("cur pod len: %d", podLen)
	}

	return nil
}

func GetIpByName(name string) string {
	ifaces, err := net.InterfaceByName(name)
	if err != nil {
		return ""
	}
	addrs, err := ifaces.Addrs()
	for _, addr := range addrs {
		ipaddr := strings.Split(addr.String(), "/")[0]
		ip := net.ParseIP(ipaddr)
		return ip.To4().String()
	}
	return ""
}
