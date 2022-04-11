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
	"github.com/yosssi/gohtml"
)

func init() {
	RootCmd.AddCommand(runCommand)
}

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "Adds SRI hashes to all <script> tags in any HTML file in the tree from the working directory",
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		paths, err := filer.Dir(pwd)
		if err != nil {
			panic(err)
		}

		if len(paths) == 0 {
			fmt.Println("No HTML files found. Exiting")
			return
		}

		for _, path := range paths {
			f, err := filer.Read(path)
			if err != nil {
				panic(err)
			}

			var ignoredHosts []string
			flag := cmd.Flag("origin")
			if flag != nil {
				ignoredHosts = []string{flag.Value.String()}
			}

			newMarkup, err := addSRIs(f, ignoredHosts)
			if err != nil {
				panic(err)
			}

			err = filer.Write(newMarkup, path)
			if err != nil {
				panic(err)
			}
		}
	},
}

func addSRIs(markup string, ignoredHosts []string) (string, error) {
	urls, err := extractor.ExtractURLS(strings.NewReader(markup))
	if err != nil {
		return "", err
	}

	f := scriptfetcher.New(ignoredHosts)
	html := markup
	integrities := map[string]string{}
	for _, u := range urls.Scripts {
		script, err := f.Fetch(u)

		if script == scriptfetcher.SKIPPED {
			continue
		}

		if err != nil {
			return "", err
		}
		h := hash.Hash([]byte(script), []crypto.Hash{crypto.SHA256, crypto.SHA384, crypto.SHA512})

		integrity := fmt.Sprintf("%v %v %v", h[crypto.SHA256], h[crypto.SHA384], h[crypto.SHA512])
		integrities[u] = integrity

		html, err = injector.Inject(html, u, integrity, "script")
		if err != nil {
			return "", err
		}
	}

	for _, u := range urls.Links {
		script, err := f.Fetch(u)

		if script == scriptfetcher.SKIPPED {
			continue
		}

		if err != nil {
			return "", err
		}
		h := hash.Hash([]byte(script), []crypto.Hash{crypto.SHA256, crypto.SHA384, crypto.SHA512})

		integrity := fmt.Sprintf("%v %v %v", h[crypto.SHA256], h[crypto.SHA384], h[crypto.SHA512])
		integrities[u] = integrity

		html, err = injector.Inject(html, u, integrity, "link")
		if err != nil {
			return "", err
		}
	}

	return gohtml.Format(html), nil
}
