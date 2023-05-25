package entity

type Dns struct {
	Kind	string	`json:"kind" yaml:"kind"`
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec     DnsSpec    `json:"spec" yaml:"spec"`
	Status   DNsStatus  `json:"status,omitempty" yaml:"status,omitempty"`
}

type DnsSpec struct {

	Host string	`json:"host" yaml:"host"`

	Paths []PathMapping	`json:"paths" yaml:"paths"`

}

type PathMapping struct {
	Path string		`json:"path" yaml:"path"`
	
	ServiceName string `json:"serviceName" yaml:"serviceName"`

	ServicePort uint16 `json:"servicePort" yaml:"servicePort"`

}

type DNsStatus struct {
	Phase string `json:"phase,omitempty" yaml:"phase,omitempty"`
}