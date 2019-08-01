package main

import (
	sm "gophercises/sitemap/sitemap"
	"log"
	"os"
)

func main() {
	xml, err := sm.GenSitemap("https://www.calhoun.io")
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(xml)
}
