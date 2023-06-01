package cmd

import (
	"fmt"
	"os"

	"minik8s/configs"
	pb "minik8s/pkg/proto"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var rootCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "kubectl is to control minik8s cluster",
	Long: `kubectl is to control minik8s cluster
		   for example:
		   kubectl get pod [podname] [-n namespace]  	get pod info of namespace
		   kubectl delete deployment [deployment name] [-n namespace] 	delete deployment
		   kubectl describe pod [pod name] [-n namespace]  get detailed info of pod`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kubectl called %s, %v", cmd.Name(), args)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(addCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewClient() pb.ApiServerKubectlServiceClient {
	conn, err := grpc.Dial(configs.ApiServerUrl+configs.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("err when connect to apiserver")
		return nil
	}
	return pb.NewApiServerKubectlServiceClient(conn)
}
