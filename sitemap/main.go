package main

import (
	"encoding/xml"
	"flag"
	"io"
	"log"
	"net/url"
	"os"
	"sort"
)

func main() {
	var maxDepth int
	flag.IntVar(&maxDepth, "depth", 3, "Number of links from the root URL to travel")
	var urlString string
	flag.StringVar(&urlString, "url", "", "URL to generate a sitemap for")

	flag.Parse()

	rootURL, err := url.Parse(urlString)
	if err != nil {
		log.Fatalf("Couldn't parse URL: '%s'\n%s", urlString, err)
	}

	sitemap := make(map[string]bool)
	sitemap[rootURL.String()] = true

	err = sitemap.Crawl(rootURL.String(), rootURL.String(), sitemap, 0, maxDepth)
	if err != nil {
		log.Fatalf("Couldn't crawl URL: '%s'\n%s", rootURL, err)
	}

	urls := make([]URL, len(sitemap))
	i := 0
	for link := range sitemap {
		urls[i].Loc = link
		i++
	}

	sort.Sort(&urlSorter{urls})

	f, err := os.Create("sitemap.xml")
	if err != nil {
		log.Fatalln("Couldn't create file: 'sitemap.xml'")
	}
	defer f.Close()

	io.WriteString(f, xml.Header)
	io.WriteString(f, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`+"\n")

	e := xml.NewEncoder(f)
	e.Indent("  ", "  ")
	e.Encode(urls)

	io.WriteString(f, "\n"+`</urlset>`+"\n")
}
