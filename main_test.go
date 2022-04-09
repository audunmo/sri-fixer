package main

import (
	"crypto"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
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

// TestExtractHashAndInject runs the primary test. It will perform the following flow:
// 1. Read all script srcs from the markup
// 2. Download all the remote scripts
// 3. Hash all the scripts
// 4. Inject integrity hashes into the markup
func TestExtractHashAndInject(t *testing.T) {
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
  integrities := map[string]string{}
  for _, u := range urls {
    script, err := f.Fetch(u)
    if err != nil {
      t.Fatal(err)
    }
    h := hash.Hash([]byte(script), []crypto.Hash{crypto.SHA256, crypto.SHA384, crypto.SHA512})

    integrity := fmt.Sprintf("%v %v %v", h[crypto.SHA256], h[crypto.SHA384], h[crypto.SHA512])
    integrities[u] = integrity

    html, err = injector.Inject(html, u, integrity)
    if err != nil {
      t.Fatal(err)
    }
  }

  newDoc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
  if err != nil {
    t.Fatal(err)
  }

  foundAndVerified := map[string]bool{}

  for url := range integrities {
    foundAndVerified[url] = false
  }

  newDoc.Find("script").Each(func(n int, s *goquery.Selection) {
    src, _ := s.Attr("src")
    expectedHash := integrities[src]
    actualHash, _ := s.Attr("integrity")

    if expectedHash == actualHash {
      foundAndVerified[src] = true
    }
  })

  for url, verified := range foundAndVerified {
    if !verified {
      t.Fatalf("was unable to find or verify the integrity hash for script %v", url)
    }
  }
}
