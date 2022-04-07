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
  exp := []string{
    "https://lmao.org/v1.3.3.7/script.js",
    "https://lmao.org/v1.3.3.7/lmao.js",
    "https://lmao.org/script.js",
  }

  sort.Strings(res)
  sort.Strings(exp)
  if !reflect.DeepEqual(res, exp) {
    t.Fatalf("Found incorrect script urls. Expected %v, got %v", exp, res)
  }
}
