package injector

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
