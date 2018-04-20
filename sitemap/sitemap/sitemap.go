package sitemap

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"

	"github.com/msiadak/gophercises/link"
)

type URL struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

type urlSorter struct {
	URLs []URL
}

func (us *urlSorter) Len() int {
	return len(us.URLs)
}

func (us *urlSorter) Less(i, j int) bool {
	return us.URLs[i].Loc < us.URLs[j].Loc
}

func (us *urlSorter) Swap(i, j int) {
	us.URLs[i], us.URLs[j] = us.URLs[j], us.URLs[i]
}

// Crawl visits the links
func Crawl(domain string, path string, sitemap map[string]bool, depth int, maxDepth int) error {
	fmt.Printf("Crawling '%s'", path)
	domainURL, err := url.Parse(domain)
	if err != nil {
		return err
	}

	resp, err := http.Get(domain)
	if err != nil {
		return err
	}

	links, err := link.ExtractLinks(resp.Body)
	if err != nil {
		return err
	}

	toCrawl := make([]string, 0, len(links))

	for _, link := range links {
		linkURL, err := domainURL.Parse(link.HREF)
		if err != nil {
			return err
		}

		linkURL.RawQuery = ""
		linkURL.Fragment = ""

		if _, ok := sitemap[linkURL.String()]; !ok && domainURL.Hostname() == linkURL.Hostname() && depth <= maxDepth {
			sitemap[linkURL.String()] = true
			if depth+1 <= maxDepth {
				toCrawl = append(toCrawl, linkURL.String())
			}
		}
	}

	for _, u := range toCrawl {
		err := Crawl(domain, u, sitemap, depth+1, maxDepth)
		if err != nil {
			return err
		}
	}

	return nil
}
