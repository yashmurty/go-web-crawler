package mock

import (
	"fmt"

	"github.com/yashmurty/go-web-crawler/core"
)

// Fetcher ...
type Fetcher struct {
	FetchFnCalled int
	MockedFetcher map[string]*mockedFetchResult
}

var _ core.Fetcher = &Fetcher{}

type mockedFetchResult struct {
	urls []string
}

// Fetch ...
func (f *Fetcher) Fetch(url string) ([]string, error) {
	f.FetchFnCalled++
	f.MockedFetcher = fetcher

	if res, ok := fetcher[url]; ok {
		return res.urls, nil
	}
	return nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated mockedFetcher.
var fetcher = map[string]*mockedFetchResult{
	"https://golang.org/": &mockedFetchResult{
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &mockedFetchResult{
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &mockedFetchResult{
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &mockedFetchResult{
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
