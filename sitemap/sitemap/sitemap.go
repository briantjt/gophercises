package sitemap

import (
	"encoding/xml"
	lp "gophercises/link/linkparser"
	"log"
	"net/http"
	"net/url"
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

func CrawlSite(u string) []string {
	res, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	parsedURL, err := url.Parse(res.Request.URL.String())
	hostname := parsedURL.Host
	if err != nil {
		log.Fatal(err)
	}
	visited := make(map[string]bool)
	mutex := &sync.Mutex{}
	links := lp.GetLinks(res.Body)
	var wg sync.WaitGroup
	for {
		mutex.Lock()
		if len(links) == 0 {
			mutex.Unlock()
			break
		}
		link := links[0].Href
		links = links[1:]
		mutex.Unlock()
		wg.Add(1)
		go func() {
			defer wg.Done()
			parsedLink, err := url.Parse(link)
			if err != nil {
				log.Fatal(err)
			}
			// Ignore anchor tags
			if parsedLink.Fragment != "" {
				return
			}
			if parsedLink.IsAbs() {
				if !isSameHost(parsedLink, hostname) {
					return
				}
			} else {
				parsedLink.Host = hostname
				parsedLink.Scheme = "https"
				link = parsedLink.String()
			}
			mutex.Lock()
			if visited[link] {
				mutex.Unlock()
				return
			}
			visited[link] = true
			res, err := http.Get(link)
			if err != nil {
				return
			}
			links = append(links, lp.GetLinks(res.Body)...)
			mutex.Unlock()
		}()
	}
	wg.Wait()
	listOfURLs := make([]string, len(visited))
	i := 0
	for key := range visited {
		listOfURLs[i] = key
	}
	return listOfURLs
}

func GenSitemap(u string) ([]byte, error) {
	sm := sitemapFactory()
	URLs := CrawlSite(u)
	for _, u := range URLs {
		sm.URLs = append(sm.URLs, &urlxml{Loc: u})
	}
	return xml.MarshalIndent(sm, "", "  ")
}
