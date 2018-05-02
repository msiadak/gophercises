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
	flag.IntVar(&workers, "workers", 4, "number of workers to fetch items with")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl, workers))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template, workers int) http.HandlerFunc {
	totalStories := numStories + workers
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("In handler")
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			fmt.Println("Failed to load top stories")
			return
		}

		idCh := make(chan int)
		go func() {
			for _, id := range ids {
				idCh <- id
			}
			close(idCh)
		}()

		storiesCh := make(chan item, workers)
		stories := make([]item, totalStories)
		done := make(chan struct{})

		go func() {
			for i := 0; i < totalStories; i++ {
				stories[i] = <-storiesCh
			}
			sort.Sort(byDescID(stories))
			for i := 0; i < workers+1; i++ {
				done <- struct{}{}
			}
			close(done)
		}()

		for i := 0; i < workers; i++ {
			go func(n int) {
				for {
					select {
					case <-done:
						return
					case id := <-idCh:
						hnItem, err := client.GetItem(id)
						if err != nil {
							return
						}

						item := parseHNItem(hnItem)
						if isStoryLink(item) {
							storiesCh <- item
						}
					}
				}
			}(i)
		}

		<-done
		fmt.Println("executing template?")
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
