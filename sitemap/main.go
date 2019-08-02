package main

import (
	"flag"
	"fmt"
	sm "gophercises/sitemap/sitemap"
	"log"
	"os"
	"sync"
)

func main() {
	var url = flag.String("url", "", "URL to generate sitemap")
	var depth = flag.Int("depth", 3, "How many recursive crawls to make")
	var crawl = flag.Bool("crawl", false, "Whether to only return a list of links")
	flag.Parse()
	visited := make(map[string]bool)
	ch := make(chan int)
	errs := make(chan error)
	m := &sync.Mutex{}
	go sm.CrawlSite(*url, *depth, visited, ch, errs, m)
	URLs := 1
	for i := 0; i < URLs; i++ {
		select {
		case additionalURLs := <-ch:
			URLs += additionalURLs
		case <-errs:
			continue
		}
	}
	siteURLs := make([]string, len(visited))
	i := 0
	for key := range visited {
		siteURLs[i] = key
		i++
	}
	if !*crawl {
		sitemap, err := sm.GenSitemap(siteURLs)
		if err != nil {
			log.Fatal(err)
		}
		os.Stdout.Write(sitemap)
	} else {
		for _, l := range siteURLs {
			fmt.Println(l)
		}
	}
}
