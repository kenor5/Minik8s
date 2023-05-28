package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	pb "minik8s/pkg/proto"
	//"minik8s/tools/log"
	"minik8s/tools/log"
	"time"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete resources cluster",
	Long: `delete resources cluster
		   for example:
		   kubectl delete deployment [deployment name] [-n namespace] 	delete deployment`,
	Run: doDelete,
}

func doDelete(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.PrintE("describe err must have 2 args")
		return
	}
	name := args[1]
	switch args[0] {
	case "po","pod","pods":
		deletePod(name)
	case "node","nodes":
		deleteNode(name)
	case "service":
		deleteService(name)
	case "function":
		deleteFunction(name)
	case "deployment","deploy":
		deleteDeployment(name)
	case "Dns", "dns":
		deleteDns(name)
	default:
		log.PrintE("delete err, no such object")
	}
}

func deletePod(name string) {
	// 通过 rpc 连接 apiserver
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := cli.DeletePod(ctx, &pb.DeletePodRequest{
		Data: []byte(name),
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Delete Pod, response ", res)
}

func deleteNode(name string) {

}

func deleteDeployment(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Print("begin delete Deployment", name)
	res, err := cli.DeleteDeployment(ctx, &pb.DeleteDeploymentRequest{
		DeploymentName: name,
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Delete Deployment, response ", res)
}

func deleteFunction(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Print("begin delete Function", name)
	_, err := cli.DeleteFunction(ctx, &pb.DeleteFunctionRequest{
		FunctionName: name,
	})

	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println("Delete Function, response ", res)    
	if err == nil {
		log.PrintS("Deleted function ", name)
	}
}

func deleteService(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	
	_, err := cli.DeleteService(ctx, &pb.DeleteServiceRequest{
		ServiceName: name,
	})

	if err != nil {
		log.PrintE(err)
	}

	if err == nil {
		log.PrintS("Delete svc ", name)
	}
}

func deleteDns(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	res, err := cli.DeleteDns(ctx, &pb.DeleteDnsRequest{
		DnsName: name,
	})

	if err != nil {
		log.PrintE(err)
	}

	fmt.Println("Delete dns, response ", res)
}
