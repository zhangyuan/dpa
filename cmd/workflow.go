package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var workflowCmd = &cobra.Command{
	Use: "workflow",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workflow called")
	},
}

func init() {
	rootCmd.AddCommand(workflowCmd)
}
