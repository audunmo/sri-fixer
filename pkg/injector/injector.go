package injector

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var TagTypeError = errors.New("TagType was not script nor link")

// Inject finds the given script tag, and then adds the given integrity hash
// Note, it skips
func Inject(markup, scriptUrl, integrity, tagType string) (string, error) {
	if tagType != "script" && tagType != "link" {
		return "", TagTypeError
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(markup))
	if err != nil {
		return "", err
	}

	if tagType == "script" {
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

			crossorigin, hasCrossOrigin := s.Attr("crossorigin")
			if !hasCrossOrigin || crossorigin != "anonymous" {
				fmt.Printf("\ncrossorigin either missing or not set to \"anonymous\" for %v. Will add it to markup", scriptUrl)
				s.SetAttr("crossorigin", "anonymous")
			}

			if val == scriptUrl {
				s.SetAttr("integrity", integrity)
			}
		})
	} else {
		doc.Find("link").Each(func(n int, s *goquery.Selection) {
			val, exists := s.Attr("href")
			if !exists {
				return
			}

			// We don't want to override existing integrity hashes. That should require manual intervention
			// Technically, this is a repeated check. We perform the same page when extracting links. But I'd rather overkill this
			// than have sri-fixer automatically change an existing integrity hash
			_, hasIntegrity := s.Attr("integrity")
			if hasIntegrity {
				return
			}

			crossorigin, hasCrossOrigin := s.Attr("crossorigin")
			if !hasCrossOrigin || crossorigin != "anonymous" {
				fmt.Printf("\ncrossorigin either missing or not set to \"anonymous\" for %v. Will add it to markup", scriptUrl)
				s.SetAttr("crossorigin", "anonymous")
			}

			if val == scriptUrl {
				s.SetAttr("integrity", integrity)
			}
		})
	}

	html, err := goquery.OuterHtml(doc.Find("*"))
	if err != nil {
		return "", err
	}

	return html, nil
}
