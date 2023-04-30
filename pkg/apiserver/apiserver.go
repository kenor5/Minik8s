package apiserver

import (
	"fmt"
	"minik8s/configs"
	"minik8s/tools/etcdctl"
	// "net"
// 
	clientv3 "go.etcd.io/etcd/client/v3"
	// "google.golang.org/grpc"
)

/*
	referenct to:
	https://blog.csdn.net/qq_43580193/article/details/127577709
*/

func Run() {
	// 开启etcd
	cli, err := etcdctl.Start(configs.EtcdStartPath)
	if err != nil {
		return
	}
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			fmt.Println("etcd close error")
		}
	}(cli)

	// 注册请求处理接口
	// listen, err := net.Listen("tcp", configs.GrpcPort)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return 
	// }
    
	// // 创建gRPC服务器
	// svr := grpc.NewServer()
}
