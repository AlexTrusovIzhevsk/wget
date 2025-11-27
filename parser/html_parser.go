package parser

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"
)

type Parser interface {
	ParseHTML(data []byte) (resources []string, links []string, err error)
}

type HtmlParser struct{}

func (p *HtmlParser) ParseHTML(data []byte) ([]string, []string, error) {
	document, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}
	resources, links := make([]string, 0), make([]string, 0)
	resources, links = p.parse(document, resources, links)
	return resources, links, nil
}

func (p *HtmlParser) parse(node *html.Node, resources []string, links []string) ([]string, []string) {
	if node.Type == html.ElementNode {
		switch node.Data {
		case "a":
			for _, attr := range node.Attr {
				if attr.Key == "href" && !isIgnoredScheme(attr.Val) {
					links = append(links, attr.Val)
				}
			}
		case "img":
			for _, attr := range node.Attr {
				if attr.Key == "src" && !isIgnoredScheme(attr.Val) {
					resources = append(resources, attr.Val)
				}
			}
		case "script":
			for _, attr := range node.Attr {
				if attr.Key == "src" {
					resources = append(resources, attr.Val)
				}
			}
		case "link":
			for _, attr := range node.Attr {
				if attr.Key == "href" && !isIgnoredScheme(attr.Val) {
					resources = append(resources, attr.Val)
				}
			}
		case "iframe":
			for _, attr := range node.Attr {
				if attr.Key == "src" && !isIgnoredScheme(attr.Val) {
					resources = append(resources, attr.Val)
				}
			}
		case "video", "audio":
			for _, attr := range node.Attr {
				if attr.Key == "src" && !isIgnoredScheme(attr.Val) {
					resources = append(resources, attr.Val)
				}
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		resources, links = p.parse(child, resources, links)
	}
	return resources, links
}

func isIgnoredScheme(value string) bool {
	return strings.HasPrefix(value, "javascript:") ||
		strings.HasPrefix(value, "mailto:") ||
		strings.HasPrefix(value, "tel:") ||
		strings.HasPrefix(value, "data:")
}
