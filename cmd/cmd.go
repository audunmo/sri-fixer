package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "sri-fixer",
	Short: "SRI Fixer automatically calculates and adds Subresource Integrity hashes to script tags in your HTML",
	Long: `
SRI Fixer automatically calculates and adds Subresource Integrity hashes to script tags in your HTML

It does this by recursively searching the working directory tree for HTML files, and downloading all URLs from third-parties found in <script src=...>.
After that, SRI Fixer will hash the downloaded javascript with sha256, sha348, and sha512.

SRI Fixer updates the files in-place, adding "integrity='sha256 sha348 sha512'" and "crossorigin='anonymous'"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
