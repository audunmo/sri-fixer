package scriptfetcher

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type FetcherConfig struct {
  ignoreHosts []string
}

type Fetcher struct {
  config FetcherConfig
}

// The result of a Fetch call if it skipped a URL
const SKIPPED = "SKIPPED"

func New(config FetcherConfig) *Fetcher {
  return &Fetcher{
    config: config,
  }
}

func (f *Fetcher) ShouldSkip(remoteURL string) (bool, error) {
  u, err := url.Parse(remoteURL)
  if err != nil {
    return false, err
  }

  host := u.Host
  for _, ignoredHost := range f.config.ignoreHosts {
    fmt.Printf("host: %v, ignoredHost: %v", host, ignoredHost)
    if host == ignoredHost {
      return true, nil
    }
  }

  return false, nil
}

func (f *Fetcher) Fetch(scriptUrl string) (string, error) {
  skip, err := f.ShouldSkip(scriptUrl)
  if err != nil {
    return "", err
  }
 
  if skip {
    return SKIPPED, nil
  }

  r, err := http.Get(scriptUrl)
  if err != nil {
    return "", err
  }

  body, err := io.ReadAll(r.Body)
  if err != nil {
    return "", err
  }

  return string(body), nil
}
