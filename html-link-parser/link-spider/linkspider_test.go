package linkspider

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	htmlContent := `<a href="/dog"><span>Something in a span</span>Text not in a span<b>Bold text!</b></a>
<a href="https://pkg.go.dev/golang.org/x/net/html">Shekebari</a>`
	htmlReader := strings.NewReader(htmlContent)
	links, err := Parse(htmlReader)

	if err != nil {
		t.Fatalf("Failed to parse links from HTML. Error: %s", err.Error())
	}

	if len(links) != 2 {
		t.Fatalf("Could not extract the right number of links from the html page. Extracted %d links", len(links))
	}

	if links[0].Href != "/dog" {
		t.Fatalf(`The 'href' attribute of the first parsed link does not match the expected value.
			The expected values is: %s, but the actual value was %s`, "/dog", links[0].Href)
	}
	if links[0].Text != "Something in a spanText not in a spanBold text!" {
		t.Fatalf("There is a problem with parsing the text inside the link. The actual parsed text is: '%s'", links[0].Text)
	}

	if links[1].Href != "https://pkg.go.dev/golang.org/x/net/html" {
		t.Fatalf(`The 'href' attribute of the first parsed link does not match the expected value.
			The expected values is: %s, but the actual value was %s`, "https://pkg.go.dev/golang.org/x/net/html", links[1].Href)
	}
	if links[1].Text != "Shekebari" {
		t.Fatalf("There is a problem with parsing the text inside the link. The actual parsed text is: '%s'", links[1].Text)
	}
}

func TestSkipComments(t *testing.T) {
	htmlReader := strings.NewReader("<a href='/hi'>Hello<!-- Greeting --></a>")
	links, _ := Parse(htmlReader)

	if links[0].Text != "Hello" {
		t.Fatalf("The actual parsed text was: %s", links[0].Text)
	}
}