package cmd

import (
	"crypto"
	"fmt"

	"github.com/audunmo/sri-fixer/pkg/hash"
	scriptfetcher "github.com/audunmo/sri-fixer/pkg/script_fetcher"
	"github.com/spf13/cobra"
)

var (
  scriptUrl string
)

func init() {
  RootCmd.AddCommand(hashScriptCmd)
  hashScriptCmd.Flags().StringVarP(&scriptUrl, "url", "u", "", "The URL of the script to hash")
  hashScriptCmd.MarkFlagRequired("url")
}

var hashScriptCmd = &cobra.Command{
  Use: "hash-script",
  Short: "Gets a remote script, and returns the integrity hashes",
  Run: func(cmd *cobra.Command, args []string) {
    f := scriptfetcher.New([]string{})
    script, err := f.Fetch(scriptUrl)
    if err != nil {
      panic(err)
    }

    h := hash.Hash([]byte(script), []crypto.Hash{crypto.SHA256, crypto.SHA3_384, crypto.SHA512})

    fmt.Printf("%v %v %v", h[crypto.SHA256], h[crypto.SHA384], h[crypto.SHA512])
  },
}
