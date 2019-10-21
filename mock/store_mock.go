package mock

import "github.com/yashmurty/go-web-crawler/core"

// Store ...
type Store struct {
	HasCrawledFnCalled int
	MockedStore        map[string]*mockedStoreResult
}

var _ core.Store = &Store{}

var store = map[string]*mockedStoreResult{}

type mockedStoreResult struct {
	*core.URLInfo
}

// HasCrawled ...
func (s *Store) HasCrawled(url string, info *core.URLInfo) bool {
	s.HasCrawledFnCalled++
	s.MockedStore = store

	// If url already exists, return true.
	if _, ok := store[url]; ok {
		return true
	}
	// Else, set the crawl status to true but return false since it's the first time.
	store[url] = &mockedStoreResult{&core.URLInfo{CrawledStatus: true}}
	return false
}
