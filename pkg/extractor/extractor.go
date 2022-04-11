package extractor

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

type URLExtraction struct {
	Scripts []string
	Links   []string
}

func ExtractURLS(markupReader io.Reader) (URLExtraction, error) {
	doc, err := goquery.NewDocumentFromReader(markupReader)
	if err != nil {
		return URLExtraction{}, err
	}

	extraction := URLExtraction{}
	doc.Find("script").Each(func(n int, s *goquery.Selection) {
		src, srcExists := s.Attr("src")

		// If an integrity hash already exists, we want to leave the script tag alone
		_, integrityExists := s.Attr("integrity")
		if srcExists && !integrityExists {
			extraction.Scripts = append(extraction.Scripts, src)
		}
	})

	doc.Find("link").Each(func(n int, s *goquery.Selection) {
		href, hrefExists := s.Attr("href")

		// If an integrity hash already exists, we want to leave the link tag alone
		_, integrityExists := s.Attr("integrity")
		if hrefExists && !integrityExists {
			extraction.Links = append(extraction.Links, href)
		}
	})

	return extraction, nil
}
