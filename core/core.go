package core

// Fetcher fetches a url.
type Fetcher interface {
	// Fetch returns a slice of URLs found on that page.
	Fetch(url string) (urls []string, err error)
}

// Store saves the URL.
type Store interface {
	// HasCrawled checks if a URL has already been crawled or not.
	HasCrawled(url string, info *URLInfo) bool
}

type URLInfo struct {
	CrawledStatus bool
	URL           string
	ParentURL     string
	ParentURLs    []string
	ChildrenURLs  []string
	Depth         int
	ChildrenInfo  []*URLInfo
}
