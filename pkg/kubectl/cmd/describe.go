package cmd

import (
	"context"
	"encoding/json"
	"minik8s/entity"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
	"time"

	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "show detailed info of object",
	Long: `show detailed info of object
		   for example:
		   kubectl describe pod [pod name] [-n namespace]  get detailed info of pod`,
	Run: doDescribe,
}

func doDescribe(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.PrintE("describe err must have 2 args")
		return
	}
	name := args[1]
	switch args[0] {
	case "po","pod","pods":
		describePod(name)
	case "node","nodes":
		describeNode(name)
	case "service":
		describeService(name)
	case "function":
		describeFunction(name)
	case "deployment","deploy":
		describeDeployment(name)
	default:
		log.PrintE("describe err, no such object")
	}
}

func describePod(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	res, err := cli.GetPod(ctx, &pb.GetPodRequest{
		PodName: name,
	})
	if err != nil {
		log.PrintE(err)
	}

	pod := &entity.Pod{}
	err = json.Unmarshal(res.PodData[0], pod)
	if err != nil {
		log.PrintE(err)
	}
	// format output yaml
	prettyjson, err := json.MarshalIndent(pod, "", "    ")
	if err != nil {
		log.PrintE(err)
	}
	log.Print(string(prettyjson))

}

func describeNode(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	res, err := cli.GetNode(ctx, &pb.GetNodeRequest{
		NodeName: name,
	})
	if err != nil {
		log.PrintE(err)
	}

	node := &entity.Node{}
	err = json.Unmarshal(res.NodeData[0], node)
	if err != nil {
		log.PrintE(err)
	}
	// format output yaml
	prettyjson, err := json.MarshalIndent(node, "", "    ")
	if err != nil {
		log.PrintE(err)
	}
	log.Print(string(prettyjson))

}

func describeDeployment(name string) {

	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	res, err := cli.GetDeployment(ctx, &pb.GetDeploymentRequest{
		DeploymentName: name,
	})
	if err != nil {
		log.PrintE(err)
	}

	pod := &entity.Deployment{}
	err = json.Unmarshal(res.Data[0], pod)
	if err != nil {
		log.PrintE(err)
	}
	// format output yaml
	prettyjson, err := json.MarshalIndent(pod, "", "    ")
	if err != nil {
		log.PrintE(err)
	}
	log.Print(string(prettyjson))

}

func describeFunction(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	res, err := cli.GetFunction(ctx, &pb.GetFunctionRequest{
		FunctionName: name,
	})
	if err != nil {
		log.PrintE(err)
	}

	pod := &entity.Function{}
	err = json.Unmarshal(res.Data[0], pod)
	if err != nil {
		log.PrintE(err)
	}
	// format output yaml
	prettyjson, err := json.MarshalIndent(pod, "", "    ")
	if err != nil {
		log.PrintE(err)
	}
	log.Print(string(prettyjson))
}

func describeService(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	res, err := cli.GetService(ctx, &pb.GetServiceRequest{
		ServiceName: name,
	})
	if err != nil {
		log.PrintE(err)
	}

	pod := &entity.Service{}
	err = json.Unmarshal(res.Data[0], pod)
	if err != nil {
		log.PrintE(err)
	}
	// format output yaml
	prettyjson, err := json.MarshalIndent(pod, "", "    ")
	if err != nil {
		log.PrintE(err)
	}
	log.Print(string(prettyjson))
}
