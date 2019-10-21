package crawler

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yashmurty/go-web-crawler/mock"
)

func TestCrawl(t *testing.T) {
	mockedStore := &mock.Store{}
	mockedFetcher := &mock.Fetcher{}

	done := make(chan bool)
	var parentURLs []string
	var parentURL string
	go Crawl("https://golang.org/", parentURL, parentURLs, 4, mockedFetcher, mockedStore, done)
	<-done

	if len(mockedStore.MockedStore) != 5 {
		t.Fatalf("Unique urls stored %d does not match the expected value of %d",
			len(mockedStore.MockedStore), 5)
	}

	if mockedFetcher.FetchFnCalled != 5 {
		t.Fatalf("Unique urls fetched %d does not match the expected value of %d",
			mockedFetcher.FetchFnCalled, 5)
	}
}

func TestExtractHost(t *testing.T) {
	t.Run("should succeed to extract the url", func(t *testing.T) {
		scheme := "https"
		host := "www.google.com"
		url := scheme + "://" + host
		extractedScheme, extractedHost := ExtractHost(url)
		if extractedScheme != scheme {
			t.Fatalf("Scheme %s does not match extracted value %s",
				scheme, extractedScheme)
		}
		if extractedHost != host {
			t.Fatalf("Host %s does not match extracted value %s",
				host, extractedHost)
		}
	})

	t.Run("should fail due to incorrect scheme", func(t *testing.T) {
		require.Panics(t, func() {
			ExtractHost("asd://www.google.com")
		})
	})

	t.Run("should fail due to incorrect host", func(t *testing.T) {
		require.Panics(t, func() {
			ExtractHost("http://@")
		})
	})
}
