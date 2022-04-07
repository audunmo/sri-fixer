package extractor

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)


func ExtractURLS(markupReader io.Reader) ([]string, error){
  doc, err := goquery.NewDocumentFromReader(markupReader)
  if err != nil {
    return []string{}, err
  }

  var urls []string
  doc.Find("script").Each(func(n int, s *goquery.Selection) {
    src, exists := s.Attr("src")
    if exists {
      urls = append(urls, src)
    }
  })

  return urls, nil
}
