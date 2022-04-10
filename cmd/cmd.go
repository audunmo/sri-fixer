package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "sri-fixer",
	Short: "SRI fixer automatically calculates and adds Subresource Integrity hashes to script tags in your HTML",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SRI Fixer. Use the `run`subcommand to add SRI hashes to all script tags in the tree of this folder")
	},
}
