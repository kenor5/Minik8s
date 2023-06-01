package NodeController

import (
	"context"
	"encoding/json"
	"minik8s/entity"
	pb "minik8s/pkg/proto"
	"minik8s/tools/etcdctl"
	"minik8s/tools/log"
	"time"
)

func (nodecontroller *NodeController)MonitorNode(){
	for {
		log.PrintS("monitor Node")
        RunningNodes := nodecontroller.GetAllLivingNodes()
        // 发送心跳信息
		for _, RunningNode := range RunningNodes {
			ctx := context.Background()
			conn := nodecontroller.GetNodeConnByName(RunningNode.Name)
            _, err := conn.SayHello(ctx, &pb.HelloRequest{
                Name : RunningNode.Name,
			})

			// 如果出错了，则更新Node的状态
			if err != nil {
				cli, err := etcdctl.NewClient()
				if err != nil {
					log.PrintE("connect to etcd error")
				}
			
				defer cli.Close()
				out, _ := etcdctl.Get(cli, "Node/"+RunningNode.Name)
				if RunningNode.Name == "" {
					out, _ = etcdctl.GetWithPrefix(cli, "Node/")
				}
			
				for _, v := range out.Kvs {
					node := &entity.Node{}
					err := json.Unmarshal(v.Value, node)
					if err != nil {
						panic("podNew unmarshel err")
					}
			
					node.Status = entity.NodeDead
					
					nodeByte, _ := json.Marshal(node)
					etcdctl.Put(cli, "Node/"+node.Name, string(nodeByte))
				}
			}

		}
		time.Sleep(30*time.Second)
	}
}