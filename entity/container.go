package entity

type Container struct {
	Name  string `json:"name" yaml:"name"`
	Image string `json:"image" yaml:"image"`
	// 镜像拉取策略。"Always"、"Never"、"IfNotPresent" 之一。
	ImagePullPolicy string        `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	VolumeMounts    []VolumeMount `json:"volumeMounts,omitempty" yaml:"volumeMounts,omitempty"`

	Ports     []ContainerPort      `json:"ports" yaml:"ports"`
	Resources ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`

	Command []string `json:"command,omitempty" yaml:"command,omitempty"`
	Args    []string `json:"args,omitempty" yaml:"args,omitempty"`
}

type ResourceRequirements struct {
	Limit   Quantity `json:"limits,omitempty" yaml:"limits,omitempty"`
	Request Quantity `json:"requests,omitempty" yaml:"requests,omitempty"`
}

type Quantity struct {
	Memory []string `json:"memory,omitempty" yaml:"memory,omitempty"`
	Cpu    []string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
}

type ContainerPort struct {
	ContainerPort string `json:"containerPort" yaml:"containerPort"`
	Protocol      string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	HostPort      string `json:"hostPort,omitempty" yaml:"hostPort,omitempty"`
}

type VolumeMount struct {
	Name      string `json:"name" yaml:"name"`
	MountPath string `json:"mountPath" yaml:"mountPath"`
}

type Volume struct {
	// 卷的名称。必须是 DNS_LABEL 且在 Pod 内是唯一的。
	Name string `json:"name" yaml:"name"`

	// 必需。path 是要创建的文件的相对路径名称。不得使用绝对路径，也不得包含 “..” 路径。 必须用 UTF-8 进行编码。相对路径的第一项不得用 “..” 开头。
	HostPath string `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`
}

//type hostPath struct {
//	path string `json:"path,omitempty" yaml:"path,omitempty"`
//}
