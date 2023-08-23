package linkspider

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Represents a link in an HTML document
type Link struct {
	Href string
	Text string
}

// Will take an HTML document and will extract all the link elements from it.
func Parse(r io.Reader) ([]Link, error) {
	html, err := html.Parse(r)
	if err != nil {
		return nil, err	
	}

	links := make([]Link, 0, 10)
	traverseHtml(html, &links)

	return links, nil
}

// Traverses the HTML DOM and collects information about links in the document.
func traverseHtml(node * html.Node, links *[]Link) {
	if (node.Type == html.ElementNode && node.Data == "a") {
		var href string
		for _, htmlAttr := range node.Attr {
			if htmlAttr.Key == "href" {
				href = htmlAttr.Val
				break
			}
		}

		var text string
		extractLinkText(node, &text)

		*links = append(*links, Link{Href: href, Text: strings.TrimSpace(text)})
		return
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverseHtml(c, links)
	}
}

// Use this method to extract the pure text from a link - tags and comments are stripped.
// This text can be used for analysis and categorisation by AI
func extractLinkText(node *html.Node, text *string) {
	if node.Type == html.TextNode {
		*text += node.Data
		return
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractLinkText(c, text)
	}
}
