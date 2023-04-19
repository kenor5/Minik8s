package entity

/*
	reference to
	https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/workload-resources/pod-v1/#PodSpec
	https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/common-definitions/object-meta/#ObjectMeta

	``的用法： https://blog.csdn.net/zsf211/article/details/106534361
*/

// Pod -----------------------------------------------------
type Pod struct {
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec     PodSpec    `json:"spec" yaml:"spec"`
	Status   PodStatus  `json:"status,omitempty" yaml:"status,omitempty"`
}

type ObjectMeta struct {
	// Name 在命名空间内必须是唯一的。创建资源时需要
	Name string `json:"name" yaml:"name"`

	// Namespace 定义了一个值空间，其中每个名称必须唯一。“omitempty”表示可为空，默认为default命名空间
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	// 可用于组织和分类（确定范围和选择）对象的字符串键和值的映射。
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`

	// Uid 是该对象在时间和空间上的唯一值。它通常由服务器在成功创建资源时生成，并且不允许使用 PUT 操作更改。
	Uid string `json:"uid,omitempty" yaml:"uid,omitempty"`
}

type PodSpec struct {
	// 可以由属于 Pod 的容器挂载的卷列表。
	Volumes []Volume `json:"volumes" yaml:"volumes"`
	// NodeName 是将此 Pod 调度到特定节点的请求。如果为空，则交给scheduler调度
	NodeName string `json:"nodeName,omitempty" yaml:"nodeName,omitempty"`
	// 属于 Pod 的容器列表。
	Containers []Container `json:"containers" yaml:"containers"`
}

type PodStatus struct {
	// Pod 被调度到的主机的 IP 地址。如果尚未被调度，则为字段为空。
	HostIp string `json:"host_ip" yaml:"host_ip"`

	// Pod 的 Phase 是对 Pod 在其生命周期中所处位置的简单、高级摘要。phase 的取值有五种可能性：
	// Pending Running Succeeded Failed Unknown
	Phase string `json:"phase" yaml:"phase"`

	// 分配给 Pod 的 IP 地址。至少在集群内可路由。如果尚未分配则为空。
	PodIp string `json:"pod_ip" yaml:"pod_ip"`
}

type Container struct {
	Name  string `json:"name" yaml:"name"`
	Image string `json:"image" yaml:"image"`
	// 镜像拉取策略。"Always"、"Never"、"IfNotPresent" 之一。
	ImagePullPolicy string        `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	VolumeMounts    []VolumeMount `json:"volume_mounts,omitempty" yaml:"volumeMounts,omitempty"`

	Ports     []ContainerPort      `json:"ports,omitempty" yaml:"ports,omitempty"`
	Resources ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`
}

type ResourceRequirements struct {
	Limit   map[string]Quantity `json:"limit" yaml:"limit"`
	Request map[string]Quantity `json:"request,omitempty" yaml:"request,omitempty"`
}

type Quantity struct {
	Cpu    int32 `json:"cpu" yaml:"cpu"`
	Memory int32 `json:"memory" yaml:"memory"`
}

type ContainerPort struct {
	ContainerPort string `json:"containerPort" yaml:"containerPort"`
	HostPort      string `json:"hostPort,omitempty" yaml:"hostPort,omitempty"`
	Protocol      string `json:"Protocol,omitempty" yaml:"protocol,omitempty"`
	HostIP        string `json:"HostIP,omitempty" yaml:"hostIP,omitempty"`
}

type VolumeMount struct {
	Name      string `json:"name" yaml:"name"`
	MountPath string `json:"mountPath" yaml:"mountPath"`
}

type Volume struct {
	// 卷的名称。必须是 DNS_LABEL 且在 Pod 内是唯一的。
	Name string `json:"name" yaml:"name"`

	// 必需。path 是要创建的文件的相对路径名称。不得使用绝对路径，也不得包含 “..” 路径。 必须用 UTF-8 进行编码。相对路径的第一项不得用 “..” 开头。
	Path string `json:"path" yaml:"path"`
}

// PodTemplate /*------------------------------------------------------*/
type PodTemplate struct {
	Metadata ObjectMeta      `json:"metadata" yaml:"metadata"`
	Template PodTemplateSpec `json:"template" yaml:"template"`
}

type PodTemplateSpec struct {
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec     PodSpec    `json:"spec" yaml:"spec"`
}

// Deployment 在 Replicaset 的基础上，进一步封装workload对象
type Deployment struct {
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
	Replicas int32 `json:"replica,omitempty" yaml:"replica,omitempty"`
}

type DeploymentStatus struct {
	// 此部署所针对的（其标签与选择算符匹配）未终止 Pod 的总数。
	Replicas int32 `json:"replicas" yaml:"replicas"`
}

type ReplicaSet struct {
	Metadata ObjectMeta       `json:"metadata" yaml:"metadata"`
	Spec     ReplicaSetSpec   `json:"spec" yaml:"spec"`
	Status   ReplicaSetStatus `json:"status" yaml:"status"`
}

type ReplicaSetSpec struct {
	// Selector 是针对 Pod 的标签查询，应与副本计数匹配。
	Selector LabelSelector `json:"selector" yaml:"selector"`
	// Template 是描述 Pod 的一个对象，将在检测到副本不足时创建此对象。
	Template PodTemplateSpec `json:"template" yaml:"template"`
	// Replicas 是预期副本的数量。这是一个指针，用于辨别显式零和未指定的值。默认为 1。
	Replicas int32 `json:"replicas,omitempty" yaml:"replicas,omitempty"`
}

type LabelSelector struct {
	// MatchLabels 映射中的单个 {key,value} 键值对相当于 matchExpressions 的一个元素，其键字段为 key，运算符为 In，values 数组仅包含 value。
	MatchLabels map[string]string `json:"matchLabels" yaml:"matchLabels"`
}

type ReplicaSetStatus struct {
	// Replicas 是最近观测到的副本数量。
	Replicas int32 `json:"replicas" yaml:"replicas"`

	// 此副本集可用副本（至少 minReadySeconds 才能就绪）的数量
	AvailableReplicas int32 `json:"availableReplicas,omitempty" yaml:"availableReplicas,omitempty"`
}
