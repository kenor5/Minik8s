package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(getCmd)
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
	  fmt.Fprintln(os.Stderr, err)
	  os.Exit(1)
	}
}