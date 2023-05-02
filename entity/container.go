package entity

type Container struct {
	Name  string `json:"name" yaml:"name"`
	Image string `json:"image" yaml:"image"`
	// 镜像拉取策略。"Always"、"Never"、"IfNotPresent" 之一。
	ImagePullPolicy string        `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	VolumeMounts    []VolumeMount `json:"volume_mounts,omitempty" yaml:"volumeMounts,omitempty"`

	Ports     []ContainerPort      `json:"ports,omitempty" yaml:"ports,omitempty"`
	Resources ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`

	Command []string `json:"command,omitempty" yaml:"command,omitempty"`
}

type ResourceRequirements struct {
	Limit   map[string]Quantity `json:"limit" yaml:"limit"`
	Request map[string]Quantity `json:"request,omitempty" yaml:"request,omitempty"`
}

type Quantity struct {
	Cpu    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}

type ContainerPort struct {
	ContainerPort string `json:"containerPort" yaml:"containerPort"`
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
