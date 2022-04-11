package injector

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestInject(t *testing.T) {
	scriptUrl := "https://test.com/lol.js"
	integrity := "fakery-fakes"
	markup := fmt.Sprintf(`
    <html>
      <head>
        <script src="%v"></script>
      </head>
      <body>
      </body>
    </html>
  `, scriptUrl)

	injected, err := Inject(markup, scriptUrl, integrity, "script")
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(injected))
	if err != nil {
		t.Fatal(err)
	}

	var passed bool
	doc.Find("script").Each(func(n int, s *goquery.Selection) {
		source, sourceExists := s.Attr("src")
		if !sourceExists {
			t.Fatalf("Could not find any script tags in injected HTML")
		}

		if source == scriptUrl {
			foundIntegrity, exists := s.Attr("integrity")
			if !exists {
				t.Fatalf("Expected %v to have an integrity hash, but found no integrity property", scriptUrl)
			}
			if foundIntegrity != integrity {
				t.Fatalf("Expected %v to have integrity %v, but found %v", scriptUrl, integrity, foundIntegrity)
			}

			passed = true
		}
	})
	if !passed {
		t.Fatal("Could not find the script tag")
	}
}
