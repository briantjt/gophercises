package main

import (
	"flag"
	"fmt"
	lp "gophercises/link/linkparser"
	"os"
)

func main() {
	filename := flag.String("file", "", "Name of the file to parse")
	flag.Parse()
	fileReader, err := os.Open(*filename)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", lp.GetLinks(fileReader))
}
