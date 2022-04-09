package main

import (
	"crypto"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/artlovecode/sri-fixer/pkg/extractor"
	"github.com/artlovecode/sri-fixer/pkg/hash"
	"github.com/artlovecode/sri-fixer/pkg/injector"
	scriptfetcher "github.com/artlovecode/sri-fixer/pkg/script_fetcher"
)

var (
  path1 = "/scripts/1.js"
  path2 = "/scripts/2.js"
  path3 = "/scripts/3.js"
  path4 = "/scripts/4.js"
  path5 = "/scripts/5.js"
)

func createTestServer() *httptest.Server {
  ts := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    w.Header().Add(http.CanonicalHeaderKey("content-type"), "application/json")
    if r.URL.Path == path1 {
      w.Write([]byte("console.log(1)"))
    }
    if r.URL.Path == path2 {
      w.Write([]byte("console.log(2)"))
    }
    if r.URL.Path == path3 {
      w.Write([]byte("console.log(3)"))
    }
    if r.URL.Path == path4 {
      w.Write([]byte("console.log(4)"))
    }
    if r.URL.Path == path5 {
      w.Write([]byte("console.log(5)"))
    }
  }))

  return ts
}
func TestExtractAndHash(t *testing.T) {
  ts := createTestServer()
  url1 := ts.URL + path1
  url2 := ts.URL + path2
  url3 := ts.URL + path3
  url4 := ts.URL + path4
  url5 := ts.URL + path5
  markup := fmt.Sprintf(`
    <html>
      <head>
        <script src="%v"></script>
        <script src="%v"></script>
        <script src="%v"></script>
      </head>
      <body>
        <script src="%v"></script>
        <script src="%v"></script>
      </body>
    </html>
  `, url1, url2, url3, url4, url5)

  urls, err := extractor.ExtractURLS(strings.NewReader(markup))
  if err != nil {
    t.Fatal(err)
  }

  f := scriptfetcher.New([]string{})
  html := markup
  for _, u := range urls {
    script, err := f.Fetch(u)
    if err != nil {
      t.Fatal(err)
    }
    h := hash.Hash([]byte(script), []crypto.Hash{crypto.SHA256, crypto.SHA384, crypto.SHA512})

    integrity := fmt.Sprintf("%v %v %v", h[crypto.SHA256], h[crypto.SHA384], h[crypto.SHA512])

    html, err = injector.Inject(html, u, integrity)
    if err != nil {
      t.Fatal(err)
    }
  }

  // TODO check the integrity hashes show up
}
