package apiserver

import (
	"log"
	"minik8s/configs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"minik8s/pkg/apiserver/client"
	pb "minik8s/pkg/proto"
)

/**************************************************************************
************************    API Server 主结构    ****************************
***************************************************************************/
type ApiServer struct {
	conn pb.KubeletApiServerServiceClient
}

var apiServer *ApiServer

func newApiServer() *ApiServer {
	newServer := &ApiServer{}
	kubelet_url := "127.0.0.1" + configs.KubeletGrpcPort
	newServer.conn, _ = ConnectToKubelet(kubelet_url)
	return newServer
}

func ApiServerObject() *ApiServer {
	if apiServer == nil {
		apiServer = newApiServer()
	}
	return apiServer
}

func (master *ApiServer) ApplyPod(in *pb.ApplyPodRequest) (*pb.StatusResponse, error) {
	// 发送消息给Kubelet
	err := client.KubeletCreatePod(apiServer.conn, in)
	if err != nil {
		log.Fatal(err)
		return &pb.StatusResponse{Status: -1}, err
	}

	return &pb.StatusResponse{Status: 0}, err
}

func (master *ApiServer) DeletePod(in *pb.DeletePodRequest) (*pb.StatusResponse, error) {
	// 发送消息给Kubelet
	err := client.KubeletDeletePod(apiServer.conn, in)
	if err != nil {
		log.Fatal(err)
		return &pb.StatusResponse{Status: -1}, err
	}

	return &pb.StatusResponse{Status: 0}, err
}



// TODO: 修改连接逻辑，正确的逻辑应该是Kubelet注册后，ApiServer获取了Kubelet的url，由此建立连接
func ConnectToKubelet(kubelet_url string) (pb.KubeletApiServerServiceClient, error) {
	// 发送消息给Kubelet
	dial, err := grpc.Dial(kubelet_url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// defer dial.Close()
	conn := pb.NewKubeletApiServerServiceClient(dial)
	return conn, nil
}
