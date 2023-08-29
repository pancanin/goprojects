package link

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

// Unifies the format of an url to a path, relative to the root domain.
// Removes any trailing forward slashes and maps from different urls, considered home urls to a common suffix.
// Example: http://www.example.com/about -> /about
// See tests for other examples.
func (l *Link) Normalize(rootDomain string) {
	// Remove the domain
	l.Href = strings.TrimPrefix(l.Href, rootDomain)

	// Remove the ending forward slash, if it's not a home url
	if l.Href != "/" {
		l.Href = strings.TrimSuffix(l.Href, "/")
	}
}

// Whether the link is from the same domain or is an external link.
// If the href of the link begins with '/' or the full domain name of the site, then we consider it from the same domain.
func (l *Link) IsSameDomainLink(rootDomain string) bool {
	return strings.HasPrefix(l.Href, "/") || strings.HasPrefix(l.Href, rootDomain) ||
		strings.TrimSpace(l.Href) == "" || strings.TrimSpace(l.Href) == "#"
}

// Will take an HTML document and will extract all the link elements from it.
func Parse(r io.Reader) ([]Link, error) {
	html, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := make([]Link, 0)
	extractLinks(html, &links)

	return links, nil
}

// Traverses the HTML DOM and collects information about links in the document.
func extractLinks(node *html.Node, links *[]Link) {
	if node.Type == html.ElementNode && node.Data == "a" {

		// Look through the attributes for the href attribute.
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
		extractLinks(c, links)
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
