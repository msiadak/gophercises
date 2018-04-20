package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Link represents a subset of information from an HTML <a> tag.
type Link struct {
	HREF string // Value of the href attribute
	Text string // Any text in the innerHTML of the tag with HTML stripped
}

// ExtractLinks reads text from r and returns the Links found in a slice.
// It searches via depth-first search and does not parse nested <a> tags.
func ExtractLinks(r io.Reader) ([]Link, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := make([]Link, 0)

	nodeDFS(node, &links)

	return links, nil
}

func nodeDFS(n *html.Node, links *[]Link) {
	if n.Type == html.ElementNode && n.DataAtom == atom.A {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				text := make([]string, 0)
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					extractInnerText(c, &text)
				}
				link := Link{attr.Val, strings.TrimSpace(strings.Join(text, " "))}
				*links = append(*links, link)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodeDFS(c, links)
	}
}

func extractInnerText(n *html.Node, text *[]string) {
	if n.Type == html.TextNode {
		*text = append(*text, strings.TrimSpace(n.Data))
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractInnerText(c, text)
	}
}
