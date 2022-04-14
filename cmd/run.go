package cmd

import (
	"crypto"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/audunmo/sri-fixer/pkg/extractor"
	"github.com/audunmo/sri-fixer/pkg/filer"
	"github.com/audunmo/sri-fixer/pkg/hash"
	"github.com/audunmo/sri-fixer/pkg/injector"
	scriptfetcher "github.com/audunmo/sri-fixer/pkg/script_fetcher"
	"github.com/spf13/cobra"
	"github.com/yosssi/gohtml"
)

var (
  FLAGNAME_ORIGIN string = "origin"
  origin string
)
func init() {
	RootCmd.AddCommand(runCommand)
  runCommand.Flags().StringVarP(&origin, FLAGNAME_ORIGIN, "o", "", "Sets the origin property. Scripts loaded from the origin will not be hashed")
  runCommand.MarkFlagRequired(FLAGNAME_ORIGIN)
}

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "Adds SRI hashes to all <script> tags in any HTML file in the tree from the working directory",
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

    gitignoreFlag := cmd.Flag("ignore")

    var gitignorePatterns []string
    if gitignoreFlag != nil {
      gitignoreData, err := filer.Read(filepath.Join(pwd, ".gitignore"))
      if err != nil {
        panic(err)
      }
      gitignorePatterns = strings.Split(gitignoreData, " ")
    }


		paths, err := filer.Dir(pwd, gitignorePatterns)
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
			if origin != "" {
				ignoredHosts = []string{origin}
			}

      fmt.Print(f)
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
  html, err = hashAndInject(urls, "script", html, f)
  if err != nil {
    return "", err
  }

  html, err = hashAndInject(urls, "link", html, f)
  if err != nil {
    return "", err
  }

	return gohtml.Format(html), nil
}

func hashAndInject (urls extractor.URLExtraction, tagType string, html string, f *scriptfetcher.Fetcher) (string, error) {
  markup := html

  var tags []string
  if tagType == "script" {
    tags = urls.Scripts
  }

  if tagType == "link" {
    tags = urls.Links
  }

	for _, u := range tags {
		script, err := f.Fetch(u)

		if script == scriptfetcher.SKIPPED {
			continue
		}

		if err != nil {
			return "", err
		}
		h := hash.Hash([]byte(script), []crypto.Hash{crypto.SHA256, crypto.SHA384, crypto.SHA512})

		integrity := fmt.Sprintf("%v %v %v", h[crypto.SHA256], h[crypto.SHA384], h[crypto.SHA512])

    markup, err = injector.Inject(markup, u, integrity, tagType)
		if err != nil {
			return "", err
		}
	}
  return markup, nil
}
