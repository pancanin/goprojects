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
	"example.com/link-spider/link"
	"example.com/link-spider/queue"
	"example.com/link-spider/treeformat"
)

func main() {
	// 1. Get a root url from a flag
	url := flag.String("url", "http://localhost:8080", "root domain of the site")
	index := flag.String("index", "/index.html", "the root of the sitemap")
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
	urlsQueue = queue.Enqueue[string](urlsQueue, *index)

	// until the queue is empty start a loop
	for len(urlsQueue) > 0 {
		currentUrl := queue.Top[string](urlsQueue)
		urlsQueue = queue.Pop[string](urlsQueue)

		if _, ok := visited[currentUrl]; ok {
			continue
		}
		visited[currentUrl] = true

		urlLinks[currentUrl] = make([]string, 0)

		// The url is ready for a request
		fullUrl := rootUrl + currentUrl
		resp, err := http.Get(fullUrl)

		if err != nil {
			printErr("Could not GET page from "+currentUrl, err)
			continue
		}
		defer resp.Body.Close()

		// Page received. Parsing html to get the links from the page.
		links, err := link.Parse(resp.Body)

		if err != nil {
			printErr("Could not get links from web page", err)
			continue
		}

		for i, link := range links {
			if link.IsSameDomainLink(rootUrl) {
				links[i].Normalize(rootUrl)

				urlsQueue = queue.Enqueue[string](urlsQueue, links[i].Href)

				urlLinks[currentUrl] = append(urlLinks[currentUrl], links[i].Href)
			}
		}
	}

	htmlTree := treeformat.Tree{Relations: urlLinks}
	htmlTree.ConstructHtmlList(*index)
	fmt.Println(htmlTree.S)
}

func printErr(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s. Error: %s", msg, err)
}
