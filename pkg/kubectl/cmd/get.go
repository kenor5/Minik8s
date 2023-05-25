package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"minik8s/entity"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
	"time"
	"minik8s/tools/prettyprint"
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
	var name string
	if len(args) == 1 {
		name = ""
	} else {
		name = args[1]
	}
	
	switch args[0] {
	case "po","pod","pods":
		getPod(name)
	case "node","nodes":
		getNode(name)
	case "service":
		getService(name)
	case "function":
		getFunction(name)
	case "deployment","deploy":
		getDeployment(name)
	case "Dns", "dns":
		getDns(name)
	case "job":
		getJob(name)
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
			PodName: name,
		})
		if err != nil {
			log.PrintE(err)
		}

		title := []string{"Name", "Status", "IP", "Age"}

		data := [][]string{}
		for _, onePod := range res.PodData {
		// prettyprint
			pod := &entity.Pod{}
			fmt.Println("Get Pod, response ", res)
			err = json.Unmarshal(onePod, pod)
			if err != nil {
				log.PrintE(err)
			}
			// 计算age,精确到秒

			age := time.Now().Sub(pod.Status.StartTime).Round(time.Second)
			data = append(data, []string{pod.Metadata.Name, pod.Status.Phase, pod.Status.PodIp, age.String()})
		}

		prettyprint.PrettyPrint(title, data)
}

func getNode(name string) {
	// 通过 rpc 连接 apiserver
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
	fmt.Println("Get Node, response ", res)
}

func getDeployment(name string) {

}

func getFunction(name string) {

}

func getService(name string) {
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
	
	title := []string{"Name", "Type", "ClusterIP"}
	data := [][]string{}
	for _, oneService := range res.Data{
		service := &entity.Service{}
		err = json.Unmarshal(oneService, service)
		if err != nil {
			log.PrintE(err)
		}
		data = append(data, []string{service.Metadata.Name, service.Spec.Type, service.Spec.ClusterIP})
	}
	prettyprint.PrettyPrint(title, data)

}

func getDns(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	res, err := cli.GetDns(ctx, &pb.GetDnsRequest{
		DnsName: name,
	})
	if err != nil {
		log.PrintE(err)
	}
	// prettyprint
	title := []string{"Name", "Host", "Subpath", "ServiceName"}

	dns := &entity.Dns{}
	err = json.Unmarshal(res.Data, dns)
	if err != nil {
		log.PrintE(err)
	}
	data := [][]string{
		
	}
	for _, v := range dns.Spec.Paths {
		data = append(data, []string{dns.Metadata.Name, dns.Spec.Host, v.Path, v.ServiceName})
	}
	prettyprint.PrettyPrint(title, data)
	
}

func getJob(name string) {
	cli := NewClient()
    if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	res, err := cli.GetJob(ctx, &pb.GetJobRequest{
		JobName: name,
	})
	if err != nil {
		log.PrintE(err)
	}
	fmt.Println("Get Job, response ", res)
}