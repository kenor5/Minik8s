package cmd

import (
	"context"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
	"time"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add info",
	Long: `add info`,
	Args: cobra.MinimumNArgs(1),
	Run:  doAdd,
}

func doAdd(cmd *cobra.Command, args []string) {
	var name string
	if len(args) == 1 {
		name = ""
	} else {
		name = args[1]
	}
	
	switch args[0] {
		case "node":
			addNode(name)
		default:
			log.PrintE("get err, no such object")
	}    
}

func addNode(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	log.Print("begin add Node", name)
	_, err := cli.AddNode(ctx, &pb.AddNodeRequest{
		NodeName: name,
	})

	if err != nil {
		log.PrintE(err)
	}

	// log.Print("Update Function, response ", res) 
	
	if err == nil {
		log.PrintS("Add Node ", name)
	}
}