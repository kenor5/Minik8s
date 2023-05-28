package cmd

import (
	// "minik8s/tools/log"
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
// TODO

}
