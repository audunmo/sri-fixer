package extractor

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestFindsURLS(t *testing.T) {
	data := `
    <html>
      <head>
        <script src="https://lmao.org/script.js"></script>
				<link href="https://lmao.org/v1.3.3.7/styles.css"></script>
      </head>
      <body>
        <script src="https://lmao.org/v1.3.3.7/script.js"></script>
        <script src="https://lmao.org/v1.3.3.7/lmao.js"></script>
        <script integrity="laaa" src="https://lmao.org/v1.3.3.7/lmao.js"></script>
      </body>
    </html>
  `
	res, err := ExtractURLS(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	expLinks := []string{
		"https://lmao.org/v1.3.3.7/styles.css",
	}
	expScripts := []string{
		"https://lmao.org/v1.3.3.7/script.js",
		"https://lmao.org/v1.3.3.7/lmao.js",
		"https://lmao.org/script.js",
	}

	sort.Strings(res.Links)
	sort.Strings(res.Scripts)
	sort.Strings(expLinks)
	sort.Strings(expScripts)
	if !reflect.DeepEqual(expScripts, res.Scripts) {
		t.Fatalf(`
Found incorrect script urls.
		Expected %v,
		got %v`,
			expScripts,
			res.Scripts)
	}

	if !reflect.DeepEqual(expLinks, res.Links) {
		t.Fatalf(`
Found incorrect script urls.
		Expected %v,
		got %v`,
			expLinks,
			res.Links)
	}
}
