package e2e

import (
	"fmt"
	"log"
	"os"
	"testing"

	crawler "github.com/yashmurty/go-web-crawler"
	"github.com/yashmurty/go-web-crawler/fetcher"
	"github.com/yashmurty/go-web-crawler/store"
)

func TestCrawler(t *testing.T) {
	InputURL := "https://yashmurty.com"
	MaxDepth := 2

	Scheme, Host := crawler.ExtractHost(InputURL)
	fmt.Println("----- Initiating web crawl for Host : ", Host)

	mapStore := &store.MemoryMapStore{
		SitemapFileName: "e2e-sitemap.txt",
	}
	mapStore.Init()

	fetcher := fetcher.URLFetcher{
		Scheme: Scheme,
		Host:   Host,
	}

	done := make(chan bool)
	var parentURLs []string
	var parentURL string
	go crawler.Crawl(InputURL, parentURL, parentURLs, MaxDepth, fetcher, mapStore, done)
	<-done

	fmt.Println("----- Finished crawling")

	mapStore.SaveToFile(InputURL)

	file, err := os.Open("e2e-sitemap.txt")
	if err != nil {
		log.Fatal(err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fi.Size())
	if fi.Size() == 0 {
		t.Fatal("Sitemap file size is zero after crawling")
	}
}
