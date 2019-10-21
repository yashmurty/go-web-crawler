package crawler

import (
	"fmt"
	"net/url"

	"github.com/yashmurty/go-web-crawler/core"
)

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url, parentURL string, parentURLs []string, depth int, fetcher core.Fetcher,
	store core.Store, done chan bool) {
	// Report to caller that we're finished.
	if done != nil {
		defer func() { done <- true }()
	}

	if depth <= 0 {
		return
	}
	// Don't fetch the same URL twice.
	info := &core.URLInfo{
		CrawledStatus: true,
		URL:           url,
		Depth:         depth,
		ParentURL:     parentURL,
		ParentURLs:    parentURLs,
		ChildrenInfo:  make([]*core.URLInfo, 0),
	}
	if store.HasCrawled(url, info) {
		return
	}

	urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Depth remaining : %d | Successfully fetched url : %s\n", depth, url)

	// Wait for the children crawls.
	childrenDone := make(chan bool, 1)
	parentURLs = append(parentURLs, url)

	for _, u := range urls {
		go Crawl(u, url, parentURLs, depth-1, fetcher, store, childrenDone)
	}
	for i := 0; i < len(urls); i++ {
		<-childrenDone
	}

	return
}

func ExtractHost(inputURL string) (string, string) {
	// Parse the URL and panic if there are errors.
	u, err := url.Parse(inputURL)
	if err != nil {
		panic("could not parse url. Please enter a valid url.")
	}
	if !(u.Scheme == "http" || u.Scheme == "https") {
		panic("Missing http(s) scheme")
	}
	if u.Host == "" {
		panic("Missing host")
	}

	return u.Scheme, u.Host
}
