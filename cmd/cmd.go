package cmd

import (
	"crypto"
	"fmt"
	"os"
	"strings"

	"github.com/artlovecode/sri-fixer/pkg/extractor"
	"github.com/artlovecode/sri-fixer/pkg/filer"
	"github.com/artlovecode/sri-fixer/pkg/hash"
	"github.com/artlovecode/sri-fixer/pkg/injector"
	scriptfetcher "github.com/artlovecode/sri-fixer/pkg/script_fetcher"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
  Use: "sri-fixer",
  Short: "SRI fixer automatically calculates and adds Subresource Integrity hashes to script tags in your HTML",
  Long: "Lorem ipsum",
  Run: func(cmd *cobra.Command, args []string) {
    pwd, err := os.Getwd()
    if err != nil {
      fmt.Println(err)
      panic(err)
    }

    paths, err := filer.Dir(pwd)
    if err != nil {
      fmt.Println(err)
      panic(err)
    }

    for _, path := range paths {
      f, err := filer.Read(path)
      if err != nil {
        fmt.Println(err)
        panic(err)
      }

      newMarkup, err := InjectSRIs(f)
      if err != nil {
        fmt.Println(err)
        panic(err)
      }

      fmt.Println(newMarkup)
    }
  },
}

func InjectSRIs(markup string) (string, error) {
  urls, err := extractor.ExtractURLS(strings.NewReader(markup))
  if err != nil {
    return "", err
  }

  f := scriptfetcher.New([]string{})
  html := markup
  integrities := map[string]string{}
  for _, u := range urls {
    script, err := f.Fetch(u)
    if err != nil {
      return "", err
    }
    h := hash.Hash([]byte(script), []crypto.Hash{crypto.SHA256, crypto.SHA384, crypto.SHA512})

    integrity := fmt.Sprintf("%v %v %v", h[crypto.SHA256], h[crypto.SHA384], h[crypto.SHA512])
    integrities[u] = integrity

    html, err = injector.Inject(html, u, integrity)
    if err != nil {
      return "", err
    }
  }

  return html, nil
}
