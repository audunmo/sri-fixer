package scriptfetcher

import (
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

func New(ignoredHosts []string) *Fetcher {
  var hosts []string
  for _, v := range ignoredHosts{
    u, err := formatHost(v)
    if err != nil {
      panic(err)
    }
    hosts = append(hosts, u)
  }
  return &Fetcher{
    config: FetcherConfig{
      ignoreHosts: hosts,
    },
  }
}

func formatHost(URL string) (string, error) {
  u, err := url.Parse(URL)
  if err != nil {
    return "", err
  }
  return u.Host, nil
}

func (f *Fetcher) ShouldSkip(remoteURL string) (bool, error) {
  host, err := formatHost(remoteURL)
  if err != nil {
    return false, err
  }

  if string(remoteURL[0]) == "/" {
    return true, nil
  }

  for _, ignoredHost := range f.config.ignoreHosts {
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
