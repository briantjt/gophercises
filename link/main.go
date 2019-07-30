package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func getText(node *html.Node, visited map[*html.Node]bool) string {
	var text string
	if node.Type == html.TextNode {
		text = node.Data
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		text += getText(child, visited)
		visited[child] = true
	}
	return strings.Trim(text, "\n")
}

func main() {
	filename := flag.String("file", "", "Name of the file to parse")
	flag.Parse()
	fileReader, err := os.Open(*filename)
	var nodes []*html.Node
	if err != nil {
		fmt.Println(err)
	}

	firstNode, err := html.Parse(fileReader)
	var links []Link
	visited := make(map[*html.Node]bool)
	if firstNode.FirstChild != nil {
		nodes = append(nodes, firstNode)
	}
	if err != nil {
		fmt.Println(err)
	}
	for len(nodes) > 0 {
		currentNode := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		if visited[currentNode] {
			continue
		}
		visited[currentNode] = true
		if currentNode.Type == html.ElementNode && currentNode.Data == "a" {
			link := Link{}
			for _, a := range currentNode.Attr {
				if a.Key == "href" {
					link.Href = a.Val
				}
			}
			link.Text = getText(currentNode, visited)
			links = append(links, link)
		}
		for child := currentNode.FirstChild; child != nil; child = child.NextSibling {
			nodes = append(nodes, child)
		}

	}
	fmt.Println(links)
}
