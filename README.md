# Go Web Crawler

A simple web crawler.

Given a URL, it outputs a simple textual `sitemap`, showing the links between pages. 

The crawler is limited to one `subdomain`. 
* When you start with *https://google.com/*, it would crawl all pages within `google.com`, but does not follow external links, for example to `facebook.com` or `mail.google.com`.

### Go Packages

```bash
/go-web-crawler/           # Contains web crawler logic.
/go-web-crawler/core       # Contains all interface definitions (Store, Fetcher)
/go-web-crawler/store      # Package store protects all shared data within a manager goroutine that accepts commands using a channel.
/go-web-crawler/fetcher    # Package fetcher fetches a url and returns children urls.
/go-web-crawler/mock       # Mock implementations of all interfaces defined in /go-web-crawler/core.
/go-web-crawler/e2e        # Our end-to-end tests.
```

### Guideline

- _Dependency Injection_ (DI): To keep things modular, easy to test, and to separate concerns, we want to be coding to interfaces, and injecting implementations of those interfaces into our `crawler` as part of our bootstrap. Using DI we're able to mock all our dependencies and unit test each layer in isolation. It's also trivial to swap out one implementation for another whenever we want.

### Usage

To run the crawler: 
```sh
go run cmd/crawler/main.go --inputURL="https://google.com" --max-depth=3 --output-file="sitemap.txt"
```
- **inputURL**: URL to be crawled. Defaults to `https://google.com`.
- **output-file**: Output file name for sitemap text file. Defaults to `sitemap.txt`.
- **max-depth**: Max depth to be crawled. Defaults to `3`.

### Tests

This section outlines how to run tests.  
We use `go modules` for dependency management. [Using Go Modules](https://blog.golang.org/using-go-modules)

> As of Go 1.11, the go command enables the use of modules when the current directory or any parent directory has a go.mod, provided the directory is outside $GOPATH/src. (Inside $GOPATH/src, for compatibility, the go command still runs in the old GOPATH mode, even if a go.mod is found. See the go command documentation for details.) Starting in Go 1.13, module mode will be the default for all development.

As noted above, please make sure that you checkout this repository outside of `$GOPATH`. Since we are using `go modules` we do not need to explicitly download the dependencies, `go build` or `go test` would take care of it.  
We use `go test` to run our tests and show us our test coverage in the output.

- Run unit tests.

  ```bash
  bash test.sh
  ```

  This will run all the unit tests in isolation with mocked dependencies.  
  Pasting output here for convenience:

  ```bash
  Running unit tests ..
  ok  	github.com/yashmurty/go-web-crawler	0.002s	coverage: 95.0% of statements
  Done.
  ```

- Run end-to-end tests.

  ```bash
  bash test.sh -e2e
  ```

  This will run all the unit tests along with tests that require real external dependencies

### Known Limitations (due to time constraints)
- We have not limited the number of go routines that can be created at the same time.    
  We achieve it by adding a worker channel and setting a maximum limit on it. Once the worker channel is exhausted, we will not create any new go routines to fetch the urls. 

- We could not add unit tests in all layers. (Store, Fetcher is still WIP.)