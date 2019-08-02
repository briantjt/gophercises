package sitemap

import (
	"encoding/xml"
	lp "gophercises/link/linkparser"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type urlxml struct {
	Loc string `xml:"loc"`
}
type sitemap struct {
	XMLName xml.Name  `xml:"urlset"`
	Xmlns   string    `xml:"xmlns,attr"`
	URLs    []*urlxml `xml:"url"`
}

func sitemapFactory() *sitemap {
	return &sitemap{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
}

func isSameHost(parsedLink *url.URL, hostname string) bool {
	if parsedLink.Scheme == "http" || parsedLink.Scheme == "https" {
		return parsedLink.Host == hostname
	}
	return false
}

func CrawlSite(u string, depth int, visited map[string]bool, ch chan int, errs chan error, mutex *sync.Mutex) {
	res, err := http.Get(u)
	if err != nil {
		errs <- err
		return
	}
	parsedURL, err := url.Parse(res.Request.URL.String())
	res.Request.Close = true
	hostname := parsedURL.Host
	if err != nil {
		errs <- err
		return
	}
	links := lp.GetLinks(res.Body)
	res.Body.Close()
	mutex.Lock()
	visited[parsedURL.String()] = true
	mutex.Unlock()
	newURLs := 0
	for _, link := range links {
		var fullURL string
		parsedLink, err := url.Parse(strings.TrimSpace(link.Href))
		if err != nil {
			errs <- err
			return
		}
		// Ignore anchor tags
		if parsedLink.Fragment != "" || strings.Contains(parsedLink.Path, ".") {
			continue
		}
		if parsedLink.IsAbs() {
			if !isSameHost(parsedLink, hostname) {
				continue
			}
		} else {
			parsedLink.Host = hostname
			parsedLink.Scheme = "https"
		}
		fullURL = parsedLink.String()
		mutex.Lock()
		if depth > 0 && !visited[fullURL] {
			go CrawlSite(fullURL, depth-1, visited, ch, errs, mutex)
			newURLs++
		}
		mutex.Unlock()
	}
	ch <- newURLs
}
func GenSitemap(urls []string) ([]byte, error) {
	sm := sitemapFactory()
	for _, u := range urls {
		sm.URLs = append(sm.URLs, &urlxml{Loc: u})
	}
	return xml.MarshalIndent(sm, "", "  ")
}
