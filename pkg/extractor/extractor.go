package extractor

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

func ExtractURLS(markupReader io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(markupReader)
	if err != nil {
		return []string{}, err
	}

	var urls []string
	doc.Find("script").Each(func(n int, s *goquery.Selection) {
		src, srcExists := s.Attr("src")

		// If an integrity hash already exists, we want to leave the script tag alone
		_, integrityExists := s.Attr("integrity")
		if srcExists && !integrityExists {
			urls = append(urls, src)
		}
	})

	doc.Find("link").Each(func(n int, s *goquery.Selection) {
		href, hrefExists := s.Attr("href")

		// If an integrity hash already exists, we want to leave the link tag alone
		_, integrityExists := s.Attr("integrity")
		if hrefExists && !integrityExists {
			urls = append(urls, href)
		}
	})

	return urls, nil
}
