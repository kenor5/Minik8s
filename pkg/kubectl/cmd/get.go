package cmd

import (
	"context"
	"fmt"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
	"time"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get resources info",
	Long: `get resources info
		   for example:
		   kubectl get pod [podname] [-n namespace]  	get pod info of namespace
		   kubectl get deployment [deployment name] [-n namespace] 	get deployment info`,
	Args: cobra.MinimumNArgs(1),
	Run:  doGet,
}

func doGet(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.LOG("get err must have 2 args")
		return
	}
	name := args[1]
	switch args[0] {
	case "po":
	case "pod":
	case "pods":
		getPod(name)
	case "node":
	case "nodes":
		getNode(name)
	case "service":
		getService(name)
	case "function":
		getFunction(name)
	case "deployment":
	case "deploy":
		getDeployment(name)
	}
}

func getPod(name string) {
		// 通过 rpc 连接 apiserver
		cli := NewClient()
		if cli == nil {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	
		res, err := cli.GetPod(ctx, &pb.GetPodRequest{
			Name: name,
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Get Pod, response ", res)
}

func getNode(name string) {

}

func getDeployment(name string) {

}

func getFunction(name string) {

}

func getService(name string) {

}
