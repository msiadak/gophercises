package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/msiadak/gophercises/quiet_hn/hn"
)

func main() {
	// parse flags
	var port, numStories, workers int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.IntVar(&workers, "workers", 8, "number of workers to fetch items with")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, workers, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func storyFetcherPool(client hn.Client, ids []int, numWorkers int) (<-chan item, chan<- struct{}) {
	idCh := make(chan int)
	storyCh := make(chan item)
	done := make(chan struct{})

	go func() {
		for _, id := range ids {
			select {
			case <-done:
				close(idCh)
				return
			case idCh <- id:
			}
		}
		close(idCh)
	}()

	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				id, more := <-idCh
				if more {
					hnItem, err := client.GetItem(id)
					if err != nil {
						log.Printf("Couldn't get item %d: %s\n", id, err)
						continue
					}
					item := parseHNItem(hnItem)
					if isStoryLink(item) {
						storyCh <- item
					}
				} else {
					return
				}
			}
		}()
	}

	return storyCh, done
}

func retrieveStories(n int, storyCh <-chan item, done chan<- struct{}) []item {
	stories := make([]item, n)
	for i := 0; i < n; i++ {
		stories[i] = <-storyCh
	}
	done <- struct{}{}
	close(done)
	return stories
}

func handler(numStories int, numWorkers int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client

		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}

		storyCh, done := storyFetcherPool(client, ids, numWorkers)
		stories := retrieveStories(numStories+(numStories/5), storyCh, done)
		sort.Sort(byDescID(stories))

		data := templateData{
			Stories: stories[:numStories],
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

type byDescID []item

func (items byDescID) Len() int {
	return len(items)
}

func (items byDescID) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items byDescID) Less(i, j int) bool {
	return items[i].ID > items[j].ID
}
