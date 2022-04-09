package injector

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Inject finds the given script tag, and then adds the given integrity hash
// Note, it skips 
func Inject(markup, scriptUrl, integrity string) (string, error) {
  doc, err := goquery.NewDocumentFromReader(strings.NewReader(markup))
  if err != nil {
    return "", err
  }

  doc.Find("script").Each(func(n int, s *goquery.Selection) {
    val, exists := s.Attr("src")
    if !exists {
      return
    }

    // We don't want to override existing integrity hashes. That should require manual intervention
    _, hasIntegrity := s.Attr("integrity")
    if hasIntegrity {
      return
    }

    if val == scriptUrl {
      s.SetAttr("integrity", integrity)
    }
  })

  html, err := goquery.OuterHtml(doc.Find("*"))
  if err != nil {
    return "", err
  }

  return html, nil
}
