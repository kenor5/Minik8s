package entity

// Deployment 在 Replicaset 的基础上，进一步封装workload对象
type Deployment struct {
	Kind	string	`json:"kind" yaml:"kind"`
	Metadata ObjectMeta       `json:"metadata" yaml:"metadata"`
	Spec     DeploymentSpec   `json:"spec" yaml:"spec"`
	Status   DeploymentStatus `json:"status" yaml:"status"`
}

type DeploymentSpec struct {
	// 供 Pod 所用的标签选择算符。
	Selector LabelSelector `json:"selector" yaml:"selector"`
	// Template 描述将要创建的 Pod。
	Template PodTemplateSpec `json:"template" yaml:"template"`
	// 预期 Pod 的数量。这是一个指针，用于辨别显式零和未指定的值。默认为 1。
	Replicas int32 `json:"replicas,omitempty" yaml:"replicas,omitempty"`
}

type DeploymentStatus struct {
	// 此部署所针对的（其标签与选择算符匹配）未终止 Pod 的总数。
	Replicas int32 `json:"replicas" yaml:"replicas"`
}
