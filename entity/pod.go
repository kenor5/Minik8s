package entity

import "time"

const (
	Pending string = "Pending"
	Running string = "Running"
	Succeed string = "Succeed"
	Failed  string = "Failed"
	Unknown string = "Unknown"
)

type Pod struct {
	Kind     string     `json:"kind" yaml:"kind"`
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec     PodSpec    `json:"spec" yaml:"spec"`
	Status   PodStatus  `json:"status,omitempty" yaml:"status,omitempty"`
}

type PodSpec struct {
	// 可以由属于 Pod 的容器挂载的卷列表。
	Volumes []Volume `json:"volumes" yaml:"volumes"`
	// NodeName 是 Pod属于哪个Node
	NodeName string `json:"nodeName,omitempty" yaml:"nodeName,omitempty"`
	// NodeSelector
	NodeSelector map[string]string 	`json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	// 属于 Pod 的容器列表。
	Containers []Container `json:"containers" yaml:"containers"`
}

type PodStatus struct {
	StartTime time.Time `json:"starttime,omitempty" yaml:"starttime,omitempty"`
	// Pod 被调度到的主机的 IP 地址。如果尚未被调度，则为字段为空。
	HostIp string `json:"host_ip,omitempty" yaml:"host_ip,omitempty"`

	// Pod 的 Phase 是对 Pod 在其生命周期中所处位置的简单、高级摘要。phase 的取值有五种可能性：
	// Pending Running Succeeded Failed Unknown
	Phase string `json:"phase,omitempty" yaml:"phase,omitempty"`

	// 分配给 Pod 的 IP 地址。至少在集群内可路由。如果尚未分配则为空。
	PodIp string `json:"pod_ip,omitempty" yaml:"pod_ip,omitempty"`
}
