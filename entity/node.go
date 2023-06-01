package entity

const (
	NodePending = "Pending"
	NodeLive = "Running"
	NodeDead = "Dead"
)

type Node struct {
	Name string  `json:"Name,omitempty" yaml:"Name,omitempty"`
	Labels map[string]string `json:"Labels,omitempty" yaml:"Labels,omitempty"`
	Ip string    `json:"Ip,omitempty"  yaml:"Ip,omitempty"`
	KubeletUrl string  `json:"KubeletUrl,omitempty" yaml:"KubeletUrl,omitempty"`
    Status string  `json:"Status,omitempty" yaml:"Status,omitempty"`
}