package kubelet

import (
	"context"
	"log"
	"minik8s/configs"
	pb "minik8s/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
    "minik8s/pkg/pod"
)

/**
** Kubelet作为客户端给ApiServer发请求
**/
func sayHello(c pb.ApiServerKubeletServiceClient) {
	ctx := context.Background()

	// 调用服务端 SimpleRPC 并获取响应
	reply, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Kubelet"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply.GetReply())
}

/**
** Kubelet作为服务端接受来自Api Server的请求
**/
type server struct {
	// 继承 protoc-gen-go-grpc 生成的服务端代码
    pb.UnimplementedKubeletApiServerServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
    log.Println("[Kubelet] Api Server call sayHello...")

	log.Println(in);
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

func (s *server) CreatePod(ctx context.Context, in *pb.ApplyPodRequest) (*pb.StatusResponse, err) {
	log.Println("[Kubelet] Api Server call Create Pod...")
    log.println(in)
    
	pod := &entity.Pod{}
	err := json.Unmarshal(in.Data, pod)
	if err != nil {
		fmt.Println("pod unmarshel err") 
		return &pb.StatusResponse{Status: -1}, err
	}

    err = RunPodSandBox(pod);
	if err != nil {
		fmt.Println("create pod err") 
		return &pb.StatusResponse{Status: -1}, err		
	}
	fmt.Println("[Kubelet] Create Pod Success")
	return &pb.StatusResponse{Status: 0}, err
}



/*********************************************************
**********************   主程序   *************************
**********************************************************/
func Run() {
	println("[kubelet] running...")

/**
*    TODO:这里的代码换成Kubelet启动时向APIServer注册
**/
	// 连接服务端，因为我们没有SSL证书，因此这里需要禁用安全传输
	apiserver_url := "127.0.0.1" + configs.GrpcPort
	dial, err := grpc.Dial(apiserver_url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer dial.Close()

	conn := pb.NewApiServerKubeletServiceClient(dial)
	sayHello(conn)

/**
*    Kubelet启动自己的服务端，接受来自ApiServer的消息
**/
	listen, err := net.Listen("tcp", configs.kubeletGrpcPort)
	if (err != nil) {
		fmt.Println(err);
	}

	// 创建gRPC服务器
	svr := grpc.NewServer()
	// 将实现的接口注册进 gRPC 服务器
    pb.RegisterKubeletApiServerServiceServer(svr,&server{})

	// 启动 gRPC 服务器
	err = svr.Serve(listen)
	if err != nil {
		log.Fatal(err)
		return
	}	
}
