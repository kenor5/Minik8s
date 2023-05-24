package configs


const (
	EtcdStartPath = "/usr/local/bin"
	ApiServerUrl = "192.168.163.132"
	NetInterface = "ens33"

	GrpcPort = ":5678"
	KubeletGrpcPort = ":5679"
  
	NginxContainerName = "n_test"
	NginxConfPath = "/root/go/src/minik8s/configs/nginx/default.conf"
	NginxLogPath = "/root/go/src/minik8s/configs/nginx/log"
    SlurmServerPort = ":8090"

)
