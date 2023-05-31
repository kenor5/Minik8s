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

			// [内存] service chain 跳转到具体 pod chain
			err = kp.IptableClient.ApplyPodRules(svcChainName, podChainName, podLen)
			if err != nil {
				return err
			}

			kp.ServiceManager.AddPodChain(svcChainName, podChainName)

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
		kp.ServiceManager.RemovePodChain(svc.ChainName)
	}

	kp.ServiceManager.RemoveServiceChain(serviceName)

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
