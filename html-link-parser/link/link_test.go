package link

import (
	"reflect"
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

func TestTrimTextSpace(t *testing.T) {
	htmlReader := strings.NewReader("<a href='/hi'>Hello  <span> What is up <span>      </br></a>")
	links, _ := Parse(htmlReader)

	if links[0].Text != "Hello   What is up" {
		t.Fatalf("The actual parsed text was: %s", links[0].Text)
	}
}

func TestLinkNormalise(t *testing.T) {
	rootUrl := "http://www.sample.com"
	tests := map[string]struct {
		input Link
		want Link
	}{
		"simple": {
			input: Link{ Href: "http://www.sample.com/about", Text: "About" },
			want: Link{Href: "/about", Text: "About"},
		},
		"trailing slash": {
			input: Link{ Href: "http://www.sample.com/about/", Text: "About" },
			want: Link{Href: "/about", Text: "About"},
		},
		"no root domain": {
			input: Link{ Href: "/about", Text: "About" },
			want: Link{Href: "/about", Text: "About"},
		},
		"just a slash": {
			input: Link{ Href: "/", Text: "Home" },
			want: Link{Href: "/", Text: "Home"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.input.Normalize(rootUrl)

			if !reflect.DeepEqual(tc.input, tc.want) {
				t.Fatalf("expected: %v, got %v", tc.want, tc.input)
			}
		})
	}
}

func TestIsSameDomainLink(t *testing.T) {
	rootDomain := "https://www.example.com/"
	tests := map[string]struct {
		input Link
		want bool
	}{
		"root domain/home page case": {
			want: true,
			input: Link{Href: rootDomain, Text: "Home"},
		},
		"relative path": {
			want: true,
			input: Link{Href: "/home", Text: "Home"},
		},
		"nested relative path": {
			want: true,
			input: Link{Href: "/page/1", Text: "First Page"},
		},
		"external link": {
			want: false,
			input: Link{Href: "http://www.whoosiewhatsie.com/home", Text: "My Blog"},
		},

		// Empty or '#' href should mean 'same document', so these are same domain links too.
		"empty href link": {
			want: true,
			input: Link{Href: "", Text: ""},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			isSameDomain := tc.input.IsSameDomainLink(rootDomain)

			if tc.want != isSameDomain {
				t.Fatalf("expected: %v, got: %v", tc.want, isSameDomain)
			}
		})
	}
}