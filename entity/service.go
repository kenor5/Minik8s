package entity

type Service struct{
	Kind	string	`json:"kind" yaml:"kind"`
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec     ServiceSpec    `json:"spec" yaml:"spec"`
	Status   ServiceStatus  `json:"status,omitempty" yaml:"status,omitempty"`
}

type ServiceSpec struct {
	Selector map[string]string 	`json:"selector" yaml:"selector"`
	Ports []ServicePort `json:"ports" yaml:"ports"`
	Type      string            `json:"type,omitempty" yaml:"type,omitempty"`
	ClusterIP string            `json:"clusterIP,omitempty" yaml:"clusterIP,omitempty"`
}

type ServicePort struct {
	Name       string `json:"name" yaml:"name"`
	Port       int32  `json:"port,omitempty" yaml:"port,omitempty"`
	TargetPort int32  `json:"targetPort,omitempty" yaml:"targetPort,omitempty"`
}

type ServiceStatus struct {
    ServicePods []string `json:"servicePods,omitempty" yaml:"servicePods,omitempty"`
}