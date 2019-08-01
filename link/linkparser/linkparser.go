package linkparser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// Recursively gets the text description of each node supplied
func getText(node *html.Node) string {
	var text string
	if node.Type == html.TextNode {
		text = node.Data
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		text += getText(child)
	}
	return strings.Trim(text, "\n")
}

func GetLinks(r io.Reader) []Link {

	firstNode, err := html.Parse(r)
	if err != nil {
		panic(err)
	}
	var nodes []*html.Node
	var links []Link

	// visited := make(map[*html.Node]bool)
	if firstNode.FirstChild != nil {
		nodes = append(nodes, firstNode)
	}
	for len(nodes) > 0 {
		currentNode := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		if currentNode.Type == html.ElementNode && currentNode.Data == "a" {
			link := Link{}
			for _, a := range currentNode.Attr {
				if a.Key == "href" {
					link.Href = a.Val
				}
			}
			link.Text = getText(currentNode)
			links = append(links, link)
		}
		for child := currentNode.FirstChild; child != nil; child = child.NextSibling {
			nodes = append(nodes, child)
		}

	}
	return links
}
