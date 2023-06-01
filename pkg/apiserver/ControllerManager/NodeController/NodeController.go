package NodeController

import (
	"encoding/json"
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/apiserver/scale"
	pb "minik8s/pkg/proto"
	"minik8s/tools/etcdctl"
	"minik8s/tools/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NodeController struct {
	NodeNameToConn map[string]pb.KubeletApiServerServiceClient
	NodeIPToConn   map[string]pb.KubeletApiServerServiceClient
	Count          int
}

func NewNodeController() *NodeController {
	newNodeController := &NodeController{
		NodeNameToConn: map[string]pb.KubeletApiServerServiceClient{},
		NodeIPToConn:   map[string]pb.KubeletApiServerServiceClient{},
		Count:          0,
	}
	return newNodeController
}

func (nodeController *NodeController) RegiseterNode(node *entity.Node) error {
	// 获取grpc kubelet的连接
	conn, err := ConnectToKubelet(node.KubeletUrl)
	if err != nil {
		panic("fail to connect kubelet: " + node.KubeletUrl)
	}
	nodeController.NodeNameToConn[node.Name] = conn
	nodeController.NodeIPToConn[node.Ip] = conn
	//连接etcd
	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("etcd client connetc error")
	}
	//將IP写入Promtheus配置
	var nodes []*entity.Node
	nodes = append(nodes, node)
	err = scale.GeneratePrometheusTargets(nodes)
	if err != nil {
		log.PrintfE("[RegiseterNode]Modify PromtheusConfiguration error!")
	}
	// 将node存入etcd
	nodeByte, err := json.Marshal(node)
	if err != nil {
		fmt.Println("parse pod error")
		return err
	}
	fmt.Printf("[ApiServer] Node %s with IP %s has registered\n", node.Name, node.Ip)
	etcdctl.Put(cli, "Node/"+node.Name, string(nodeByte))

	return nil
}

// RoundRobin调度策略
func (nodeController *NodeController) RoundRobin(nodeSelector map[string]string) (pb.KubeletApiServerServiceClient, string) {
	LivingNodes := nodeController.GetAllLivingNodes()

	// 选取Label
	var selectedNodes []*entity.Node
	for _, node := range LivingNodes {
		nodeLabels := node.Labels
		match := true

		// 检查Node的Label字段是否包含Pod的NodeSelector字段
		for key, value := range nodeSelector {
			if nodeLabels[key] != value {
				match = false
				break
			}
		}

		if match {
			selectedNodes = append(selectedNodes, node)
		}
	}    

	selectedNodesNum := len(selectedNodes)

	selectedNode := selectedNodes[nodeController.FetchAndAdd()%selectedNodesNum]
	return nodeController.NodeNameToConn[selectedNode.Name], selectedNode.Ip
}

func (nodeController *NodeController) FetchAndAdd() int {
	result := nodeController.Count
	nodeController.Count += 1
	return result
}

// 获取所有活跃着的Node
func (nodeController *NodeController) GetAllLivingNodes() []*entity.Node {
	LivingNodes := []*entity.Node{}

	// 从etcd中拿出所有的Pod
	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("etcd client connetc error")
	}

	out, _ := etcdctl.GetWithPrefix(cli, "Node/")

	// 筛选出符合Label且状态为Running的Pod
	for _, data := range out.Kvs {
		node := &entity.Node{}
		err = json.Unmarshal(data.Value, node)
		if err != nil {
			fmt.Println("pod unmarshal error")
		}
		//fmt.Println("get etcd", pod)

		// 判断Pod仍在运行(状态为Running)Selector和Label完全相等
		if node.Status == entity.NodeLive {
			LivingNodes = append(LivingNodes, node)
		}
	}

	return LivingNodes
}

// 将Node存入etcd中
func (nodeController *NodeController) SetNode(Node *entity.Node) (error) {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connetc error")
		return err
	}
	defer cli.Close()
	nodeByte, err := json.Marshal(Node)
	etcdctl.Put(cli, "Node/"+Node.Name, string(nodeByte))
	return nil
}

// 根据Node的名称获取conn
func (nodeController *NodeController) GetNodeConnByName(NodeName string) pb.KubeletApiServerServiceClient {
	return nodeController.NodeNameToConn[NodeName]
}

// 根据Node的名称获取conn
func (nodeController *NodeController) GetNodeConnByIP(NodeIP string) pb.KubeletApiServerServiceClient {
	return nodeController.NodeIPToConn[NodeIP]
}

// 重新获取Node Conn
func (nodeController *NodeController) RestartNodeConn() (error) {
    LivingNodes := nodeController.GetAllLivingNodes()

	for _, LivingNode := range LivingNodes {
		// 获取grpc kubelet的连接
	    conn, err := ConnectToKubelet(LivingNode.KubeletUrl)
		if err != nil {
			log.PrintW("fail to connect kubelet: " + LivingNode.KubeletUrl)
		    // 更新etcd中Node的状态
			LivingNode.Status= entity.NodeDead
			// 存入etcd中
            nodeController.SetNode(LivingNode)
			continue
		}
		// 加入内存中的map
		nodeController.NodeNameToConn[LivingNode.Name] = conn
		nodeController.NodeIPToConn[LivingNode.Ip] = conn
	}

	return nil
}

// 工具函数
func ConnectToKubelet(kubelet_url string) (pb.KubeletApiServerServiceClient, error) {
	// 发送消息给Kubelet
	dial, err := grpc.Dial(kubelet_url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.PrintE(err)
		return nil, err
	}
	// defer dial.Close()
	conn := pb.NewKubeletApiServerServiceClient(dial)
	return conn, nil
}
