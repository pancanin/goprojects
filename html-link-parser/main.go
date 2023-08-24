// Ok, I am too lazy to create separate packages for this one...I will code the sitemap builder here
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	// Still learning the package system
	// So the first part 'example.com/link-spider' is the module that I init-ed.
	// and the part after that is the folder where the package I want to import is.
	// I butchered this project, but I will try to make it good in the next one.
	linkspider "example.com/link-spider/link"
)

func main() {
	// 1. Get a root url from a flag
	url := flag.String("url", "https://www.calhoun.io/", "the url which will be the root of the sitemap.")
	flag.Parse()

	if *url == "" {
		fmt.Fprintln(os.Stderr, "url is a required parameter")
		os.Exit(1)
	}

	// The root URL might not end on '/', but for starters we will assume that it is always the 'home' page.
	rootUrl := strings.TrimSuffix(*url, "/")
	// Define a queue of strings
	urls := make([]string, 0)
	// add the home page of the site to the queue
	urls = enqueue(urls, "/")
	visited := make(map[string]bool, 0)
	
	// until the queue is empty start a loop
	for len(urls) > 0 {
		currentUrl := top(urls)
		urls = pop(urls)
		
		// Prepare the url for querying
		// TODO Create a method for this part and add unit tests
		// Remove the domain
		currentUrl = strings.TrimPrefix(currentUrl, rootUrl)
		// Remove the ending forward slash
		currentUrl = strings.TrimSuffix(currentUrl, "/")

		// Map different variants of the home url to a single form
		if currentUrl == "#" || currentUrl == "" {
			currentUrl = "/"
		}

		if _, ok := visited[currentUrl]; ok {
			continue
		}
		visited[currentUrl] = true

		// Instead, we should have a map from urls to child links
		fmt.Println(currentUrl)

		// The url is ready for querying
		fullUrl := rootUrl + currentUrl
		resp, err := http.Get(fullUrl)

		if err != nil {
			fmt.Println("Could not get page from " + currentUrl)
			continue
		}
		defer resp.Body.Close()
		links, err := linkspider.Parse(resp.Body)

		if err != nil {
			fmt.Println("Could not get links from web page")
			continue
		}

		// Filter just the same domain links
		var sameDomainLinks []linkspider.Link

		for _, link := range links {
			if strings.HasPrefix(link.Href, "/") || strings.HasPrefix(link.Href, *url) {
				sameDomainLinks = append(sameDomainLinks, link)
			}
		}

		// Print the links at the current page
		for i, link := range sameDomainLinks {
			fmt.Printf("Found link with address %s and text %s\n", link.Href, link.Text)
			sameDomainLinks[i].Href = strings.TrimPrefix(link.Href, rootUrl)
			sameDomainLinks[i].Href = strings.TrimSuffix(link.Href, "/")

			if link.Href == "#" || link.Href == "" {
				sameDomainLinks[i].Href = "/"
			}

			if _, ok := visited[link.Href]; !ok {
				urls = enqueue(urls, link.Href)
			}
		}
	}
}

type Page struct {
	Url string
	Children []*Page
}

func enqueue(queue []string, item string) []string {
	return append(queue, item)
}

func top(queue []string) string {
	return queue[0]
}

func pop(queue []string) []string {
	return queue[1:]
}

// TODO: Look at the identifiers around your code and improve them.
// TODO: Reduce abstractions to improve readability
// Use 'internal' directory if we dont want to share an implementation with the outside world