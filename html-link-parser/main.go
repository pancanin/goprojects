package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"path/filepath"

	// Still learning the package system
	// So the first part 'example.com/link-spider' is the module that I init-ed.
	// and the part after that is the folder where the package I want to import is.
	"example.com/link-spider/link"
	"example.com/link-spider/queue"
)

var urlLinks map[string][]string

func main() {
	// 1. Get a root url from a flag
	url := flag.String("url", "https://calhoun.io/", "the url which will be the root of the sitemap.")
	flag.Parse()

	if *url == "" {
		fmt.Fprintln(os.Stderr, "url is a required parameter")
		os.Exit(1)
	}

	// The root URL might not end on '/', but for starters we will assume that it is always the 'home' page.
	rootUrl := strings.TrimSuffix(*url, "/")

	urlsQueue := make([]string, 0)
	visited := make(map[string]bool, 0)
	urlLinks := make(map[string][]string, 0)

	// add the home page of the site to the queue
	urlsQueue = queue.Enqueue[string](urlsQueue, "/")
	
	//until the queue is empty start a loop
	for len(urlsQueue) > 0 {
		currentUrl := queue.Top[string](urlsQueue)
		urlsQueue = queue.Pop[string](urlsQueue)

		if _, ok := visited[currentUrl]; ok {
			continue
		}
		visited[currentUrl] = true

		urlLinks[currentUrl] = make([]string, 0)

		// The url is ready for a request
		fullUrl := filepath.Join(rootUrl, currentUrl)
		resp, err := http.Get(fullUrl)

		if err != nil {
			printErr("Could not GET page from " + currentUrl, err)
			continue
		}
		defer resp.Body.Close()

		// Page received. Parsing html to get the links from the page.
		links, err := link.Parse(resp.Body)

		if err != nil {
			printErr("Could not get links from web page", err)
			continue
		}

		// We do not need to process external links.
		var sameDomainLinks []link.Link

		for _, link := range links {
			if strings.HasPrefix(link.Href, "/") || strings.HasPrefix(link.Href, *url) {
				sameDomainLinks = append(sameDomainLinks, link)
			}
		}

		// Visit the same domain links at the current page
		for i := range sameDomainLinks {
			// Remove the domain
			sameDomainLinks[i].Href = strings.TrimPrefix(sameDomainLinks[i].Href, rootUrl)
			// Remove the ending forward slash
			sameDomainLinks[i].Href = strings.TrimSuffix(sameDomainLinks[i].Href, "/")

			// Map different variants of the home url to a single form
			if sameDomainLinks[i].Href == "#" || sameDomainLinks[i].Href == "" {
				sameDomainLinks[i].Href = "/"
			}

			urlsQueue = queue.Enqueue[string](urlsQueue, sameDomainLinks[i].Href)

			urlLinks[currentUrl] = append(urlLinks[currentUrl], sameDomainLinks[i].Href)
		}
	}
	
	var html string
	constructHtmlList("/", &html)
	fmt.Println(html)
}

type Page struct {
	Url string
	Children []*Page
}

// Operates on a package local variable 'urlLinks'
// Improve this by wrapping this function in an object perhaps
func constructHtmlList(url string, html *string) {
	*html += "<ul>"
	*html += "<li>" + url

	if _, ok := urlLinks[url]; ok && len(urlLinks[url]) > 0 {
		// there are children
		for _, link := range urlLinks[url] {
			constructHtmlList(link, html)
			*html += "</ul>"
		}
	} else {
		*html += "</li>"
	}
}

// TODO: Look at the identifiers around your code and improve them.
// TODO: Reduce abstractions to improve readability
// Use 'internal' directory if we dont want to share an implementation with the outside world

func printErr(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s. Error: %s", msg, err.Error())
}
