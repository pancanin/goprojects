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
	url := flag.String("url", "https://example.com", "the home url of the website for which sitemap will be generated.")
	flag.Parse()

	if *url == "" {
		fmt.Fprintln(os.Stderr, "url is a required parameter")
		os.Exit(1)
	}

	// 2. Visit the page and find all links in it, then visit the links and do the same
	resp, err := http.Get(*url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "There was an error requesting resource at %s. Error: %s", *url, err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	links, err := linkspider.Parse(resp.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "There was an error parsing html file from %s. Error: %s", *url, err.Error())
		os.Exit(1)
	}

	// Filter just the same domain links
	var sameDomainLinks []linkspider.Link

	for _, link := range links {
		if strings.HasPrefix(link.Href, "/") || strings.HasPrefix(link.Href, *url) {
			sameDomainLinks = append(sameDomainLinks, link)
		}
	}

	// Print the links at the current page
	for _, link := range sameDomainLinks {
		fmt.Printf("Found link with address %s and text %s\n", link.Href, link.Text)
	}

	// 3. While doing this create a tree-like structure with the path which we traced along the way.

	// 4. Check for already visited links, maybe a set/map of the visited pages is a good idea.

	// Notes: Visit only links from the same domain

	// struct Page { url: "the url of the page", children []Page }
}