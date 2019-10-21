package main

import (
	"flag"
	"fmt"

	crawler "github.com/yashmurty/go-web-crawler"

	"github.com/yashmurty/go-web-crawler/fetcher"
	"github.com/yashmurty/go-web-crawler/store"
)

var (
	InputURL   string
	OutputFile string
	MaxDepth   int
)

func main() {
	flag.StringVar(&InputURL, "inputURL", "https://google.com", "URL to be crawled")
	flag.StringVar(&OutputFile, "outputFile", "sitemap.txt", "Output file name for sitemap text file")
	flag.IntVar(&MaxDepth, "max-depth", 3, "Max depth to be crawled")
	flag.Parse()

	Scheme, Host := crawler.ExtractHost(InputURL)
	fmt.Println("----- Initiating web crawl for Host : ", Host)

	mapStore := &store.MemoryMapStore{
		SitemapFileName: OutputFile,
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

	return
}
