// Package fetcher fetches a url and returns a slice of all urls
// found on that page.
package fetcher

import (
	"fmt"
	"net/http"
	netURL "net/url"
	"strings"

	"golang.org/x/net/html"
)

// URLFetcher is Fetcher that fetches HTTP URLs.
type URLFetcher struct {
	Scheme string
	Host   string
}

// Fetch fetches a given URL.
func (f URLFetcher) Fetch(url string) (urls []string, err error) {
	fmt.Println("Fetch called for url : ", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: Failed to fetch \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // close response body when the function returns.
	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done.
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag.
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one.
			ok, url := getHrefTag(t)
			if !ok {
				continue
			}

			// Parse the URL and skip if there are errors.
			u, err := netURL.Parse(url)
			if err != nil {
				continue
			}

			// If the URL is just the path, then add the scheme and host as pre-fix.
			if u.Scheme == "" && u.Host == "" {
				u.Scheme = f.Scheme
				u.Host = f.Host
				if strings.HasPrefix(url, "/") {
					url = f.Scheme + "://" + f.Host + url
				} else {
					url = f.Scheme + "://" + f.Host + "/" + url
				}
			}
			if u.Host != f.Host {
				// Skip the URL if the host is outside of the given sub-domain.
				continue
			}

			// Make sure the URLs have a scheme.
			if u.Scheme == "" {
				url = f.Scheme + ":" + url
			}

			urls = append(urls, url)
		}
	}
}

// Helper function to pull the href attribute from a Token
func getHrefTag(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}
