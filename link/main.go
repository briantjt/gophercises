package main

import (
	"flag"
	"fmt"
	lp "gophercises/link/linkparser"
	"os"

	"golang.org/x/net/html"
)

func main() {
	filename := flag.String("file", "", "Name of the file to parse")
	flag.Parse()
	fileReader, err := os.Open(*filename)
	if err != nil {
		fmt.Println(err)
	}

	firstNode, err := html.Parse(fileReader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lp.GetLinks(firstNode))
}
