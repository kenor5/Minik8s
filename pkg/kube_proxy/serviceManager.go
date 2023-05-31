package kube_proxy

import (
	"minik8s/entity"
)

type ServiceChain struct {
	ChainName string
	Ports     *entity.ServicePort
}

type PodChain struct {
	ChainName string
}

type ServiceManager struct {
	ServiceName2ServiceChain  map[string][]*ServiceChain
	ServiceChainName2PodChain map[string][]*PodChain
	ServiceName2ClusterIp     map[string]string
}

func NewSvcManager() *ServiceManager {
	return &ServiceManager{
		ServiceName2ServiceChain:  map[string][]*ServiceChain{},
		ServiceChainName2PodChain: map[string][]*PodChain{},
		ServiceName2ClusterIp:     map[string]string{},
	}
}

func (sm *ServiceManager) GetServiceChains(svcName string) []*ServiceChain {
	return sm.ServiceName2ServiceChain[svcName]
}

func (sm *ServiceManager) GetPodChains(svcChainName string) []*PodChain {
	return sm.ServiceChainName2PodChain[svcChainName]
}

func (sm *ServiceManager) GetClusterIp(svcName string) string {
	return sm.ServiceName2ClusterIp[svcName]
}

func (sm *ServiceManager) RemoveServiceChain(svcName string) {
	delete(sm.ServiceName2ServiceChain, svcName)
}

func (sm *ServiceManager) RemovePodChain(svcChainName string) {
	delete(sm.ServiceChainName2PodChain, svcChainName)
}

func (sm *ServiceManager) AddServiceChain(svcName string, svcChainName string, port *entity.ServicePort) {
	if !sm.ExistServiceChain(svcName) {
		sm.ServiceName2ServiceChain[svcName] = []*ServiceChain{}
	}

	sm.ServiceName2ServiceChain[svcName] = append(sm.ServiceName2ServiceChain[svcName],
		&ServiceChain{ChainName: svcChainName, Ports: port})

}

func (sm *ServiceManager) AddPodChain(svcChainName string, podChainName string) {
	if !sm.ExistPodChain(svcChainName) {
		sm.ServiceChainName2PodChain[svcChainName] = []*PodChain{}
	}
	sm.ServiceChainName2PodChain[svcChainName] = append(sm.ServiceChainName2PodChain[svcChainName],
		&PodChain{ChainName: podChainName},
	)

}

func (sm *ServiceManager) AddClusterIp(svcName string, clusterIp string) {
	if !sm.ExistClusterIp(clusterIp)	{
		sm.ServiceName2ClusterIp[svcName] = clusterIp
	}

}

func (sm *ServiceManager) ExistServiceChain(svcChainName string) bool {
	_, exist := sm.ServiceName2ServiceChain[svcChainName]
	return exist
}

func (sm *ServiceManager) ExistPodChain(svcChainName string) bool {
	_, exist := sm.ServiceChainName2PodChain[svcChainName]
	return exist
}

func (sm *ServiceManager) ExistClusterIp(svcName string) bool {
	_, exist := sm.ServiceName2ClusterIp[svcName]
	return exist
}
