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
		default:
			log.PrintE("get err, no such object")
	}    
    
}

func updateFunction(name string) {
	cli := NewClient()
	if cli == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Print("begin update Function", name)
	res, err := cli.UpdateFunction(ctx, &pb.UpdateFunctionRequest{
		FunctionName: name,
	})

	if err != nil {
		log.PrintE(err)
	}

	log.Print("Update Function, response ", res)    

}