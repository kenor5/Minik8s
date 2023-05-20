package entity

const (
	NodeLive = "Living"
	NodeDead = "Dead"
)

type Node struct {
	Name string  `json:"Name"`
	Ip string    `json:"Ip"`
	KubeletUrl string  `json:"KubeletUrl"`
    Status string  `json:"Status"`
}