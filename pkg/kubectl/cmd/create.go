package cmd

import "github.com/spf13/cobra"

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "kubectl create is to create object in cluster",
	Long: `kubectl create is to create object in cluster
		   for example:
		   kubectl create -f *.yaml  	create pod`,
	Run: doCreate,
}

func doCreate(cmd *cobra.Command, args []string) {

}