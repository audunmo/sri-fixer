package scriptfetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestServer(dataToServe string) *httptest.Server {
  ts := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    w.Header().Add(http.CanonicalHeaderKey("content-type"), "application/json")
    w.Write([]byte(dataToServe))
  }))

  return ts
}

func TestShouldSkipIgnoredHost(t *testing.T) {
  fetcher := New([]string{"https://test.com"})

  skip, err := fetcher.ShouldSkip("https://test.com/script.js")
  if err != nil {
    t.Fatal(err)
  }

  if !skip {
    t.Fatal("Didn't skip URL marked for skip")
  }
}

func TestShouldSkipLocalFile(t *testing.T) {
  fetcher := New([]string{})

  skip, err := fetcher.ShouldSkip("/script.js")
  if err != nil {
    t.Fatal(err)
  }

  if !skip {
    t.Fatal("Didn't skip URL marked for skip")
  }
}


// Subdomains should be included, unless explicitly listed as skipped
func TestShouldNotSkipSubdomainOnIgnoredHost(t *testing.T) {
  fetcher := New([]string{"https://test.com"})

  skip, err := fetcher.ShouldSkip("https://sub.test.com/script.js")
  if err != nil {
    t.Fatal(err)
  }

  if skip {
    t.Fatal("Skipped subdomain not explicitly marked for skipping")
  }
}

func TestShouldRequestScript(t *testing.T) {
  script := "console.log('hello');"
  ts := createTestServer(script)

  f := New([]string{})

  resp, err := f.Fetch(ts.URL)
  if err != nil {
    t.Fatal(err)
  }

  if resp != script {
    t.Fatalf("Expected fetcher to get %v, but got %v", script, resp)
  }
}

func TestShouldReturnSkippedIfIgnored(t *testing.T) {
  script := "console.log('hello');"
  ts := createTestServer(script)

  f := New([]string{ts.URL})

  resp, err := f.Fetch(ts.URL)
  if err != nil {
    t.Fatal(err)
  }

  if resp != SKIPPED {
    t.Fatalf("Expected fetcher to get %v, but got %v", SKIPPED, resp)
  }
}
