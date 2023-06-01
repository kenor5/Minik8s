package cmd

import (
	"context"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
	"time"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update info",
	Long: `update info`,
	Args: cobra.MinimumNArgs(1),
	Run:  doUpdate,
}

func doUpdate(cmd *cobra.Command, args []string) {
	var name string
	if len(args) == 1 {
		name = ""
	} else {
		name = args[1]
	}
	
	switch args[0] {
		case "function":
			updateFunction(name)
		case "svc":
			updateSvc(name, args[2])
		case "svc2":
			updateSvc2(name, args[2])
		default:
			log.PrintE("get err, no such object")
	}    
    
}

func updateFunction(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	log.Print("begin update Function", name)
	_, err := cli.UpdateFunction(ctx, &pb.UpdateFunctionRequest{
		FunctionName: name,
	})

	if err != nil {
		log.PrintE(err)
	}

	// log.Print("Update Function, response ", res) 
	
	if err == nil {
		log.PrintS("Updated function ", name)
	}

}

// add pod to service
func updateSvc(svcName string, podName string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	log.Print("begin update svc", svcName)

	_, err := cli.UpdateSvc(ctx, &pb.UpdateSvcRequest{
		SvcName: svcName,
		PodName: podName,
	})

	if err != nil {
		log.PrintE(err)
	}

}

// del pod from service
func updateSvc2(svcName string, podName string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	log.Print("begin update svc", svcName)

	_, err := cli.UpdateSvc2(ctx, &pb.UpdateSvcRequest{
		SvcName: svcName,
		PodName: podName,
	})

	if err != nil {
		log.PrintE(err)
	}
}