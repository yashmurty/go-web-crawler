// Package store protects all shared data within a manager
// goroutine that accepts commands using a channel.
package store

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/yashmurty/go-web-crawler/core"
)

type CommandType int

const (
	GetCommand = iota
	SetCommand
	SaveToFileCommand
)

// Command will be sent across the channels.
type Command struct {
	CmdType   CommandType
	URL       string
	Info      *core.URLInfo
	ReplyChan chan *core.URLInfo
}

type MemoryMapStore struct {
	SitemapFileName string
	// URLManager is our manager for our URL store. It maintains a map that stores the URLs.
	URLManager chan<- Command
}

// Init initializes MemoryMapStore.
func (m *MemoryMapStore) Init() {
	fmt.Println("> URLManager has been initiated.")
	m.URLManager = m.startURLManager(map[string]core.URLInfo{})
}

// startURLManager starts a goroutine that serves as a manager for our
// URL store. Returns a channel that's used to send commands to the
// manager.
func (m *MemoryMapStore) startURLManager(initvals map[string]core.URLInfo) chan<- Command {
	urlMap := make(map[string]*core.URLInfo)
	for k, v := range initvals {
		urlMap[k] = &v
	}

	cmds := make(chan Command)

	go func() {
		for cmd := range cmds {
			switch cmd.CmdType {
			case GetCommand:
				if val, ok := urlMap[cmd.URL]; ok {
					cmd.ReplyChan <- val
				} else {
					cmd.ReplyChan <- &core.URLInfo{CrawledStatus: false}
				}
			case SetCommand:
				urlMap[cmd.URL] = cmd.Info
				cmd.ReplyChan <- cmd.Info
			case SaveToFileCommand:
				err := writeToFile(m.SitemapFileName, cmd.URL, urlMap)
				if err != nil {
					log.Fatal(err)
				}
				cmd.ReplyChan <- &core.URLInfo{}
			default:
				log.Fatal("unknown command type", cmd.CmdType)
			}
		}
	}()
	return cmds
}

// HasCrawled checks if a URL has already been crawled or not.
func (m *MemoryMapStore) HasCrawled(url string, info *core.URLInfo) bool {

	replyChan := make(chan *core.URLInfo)
	m.URLManager <- Command{CmdType: GetCommand, URL: url, ReplyChan: replyChan}
	reply := <-replyChan

	// If url already exists, return true.
	if reply.CrawledStatus == true {
		return true
	}
	// Else, set the crawl status to true but return false for existing status.
	m.URLManager <- Command{CmdType: SetCommand, URL: url, Info: info, ReplyChan: replyChan}
	_ = <-replyChan

	return false
}

// SaveToFile will save the map to a text file.
func (m *MemoryMapStore) SaveToFile(inputURL string) bool {
	replyChan := make(chan *core.URLInfo)
	m.URLManager <- Command{CmdType: SaveToFileCommand, URL: inputURL, ReplyChan: replyChan}
	<-replyChan

	return true
}

func writeToFile(filename string, inputURL string, data map[string]*core.URLInfo) error {
	fmt.Println("Saving resutls to file :", filename)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	sitemapHeader := "Sitemap for inputURL : " + inputURL + "\n"
	_, err = io.WriteString(file, sitemapHeader)
	if err != nil {
		return err
	}

	for u, uinfo := range data {
		if data[uinfo.ParentURL] != nil {
			data[uinfo.ParentURL].ChildrenInfo = append(data[uinfo.ParentURL].ChildrenInfo, data[u])
		}
	}

	err = writeEntryRecursive(file, inputURL, data)
	if err != nil {
		return err
	}

	return file.Sync()
}

func writeEntryRecursive(file *os.File, url string, data map[string]*core.URLInfo) error {
	if len(data) == 0 {
		return nil
	}
	leftPadding := strings.Repeat("--", len(data[url].ParentURLs)+1)
	sitemapEntry := leftPadding + "> " + url + "\n"
	_, err := io.WriteString(file, sitemapEntry)
	if err != nil {
		return err
	}

	for _, v := range data[url].ChildrenInfo {
		writeEntryRecursive(file, v.URL, data)
	}

	return nil
}
