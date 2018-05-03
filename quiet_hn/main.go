package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

	c := newCache(hn.Client{}, workers, numStories)
	c.Update()
	c.UpdateEvery(time.Minute * 10)

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(c, tpl))

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

func handler(c *cache, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		data := templateData{
			Stories: c.Get()[:30],
			Time:    time.Now().Sub(start),
		}
		err := tpl.Execute(w, data)
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

type byIDSlice struct {
	ids   []int
	items *[]item
}

func (s byIDSlice) Len() int {
	return len(*s.items)
}

func (s byIDSlice) Swap(i, j int) {
	(*s.items)[i], (*s.items)[j] = (*s.items)[j], (*s.items)[i]
}

func (s byIDSlice) Less(i, j int) bool {
	iIDPos, jIDPos := -1, -1
	for k := 0; iIDPos != -1 && jIDPos != -1; k++ {
		switch s.ids[k] {
		case (*s.items)[i].ID:
			iIDPos = k
		case (*s.items)[j].ID:
			jIDPos = k
		}
	}
	return iIDPos < jIDPos
}

type cache struct {
	client      hn.Client
	workers     int
	stories     int
	lastUpdated time.Time
	mutex       *sync.Mutex
	items       []item
}

func newCache(client hn.Client, workers, stories int) *cache {
	return &cache{
		client:  client,
		workers: workers,
		stories: stories,
		mutex:   &sync.Mutex{},
		items:   make([]item, stories),
	}
}

func (c *cache) Update() error {
	fmt.Println("Updating")
	ids, err := c.client.TopItems()
	if err != nil {
		return fmt.Errorf("Failed to load top stories: %s", err)
	}

	fmt.Println("5 ids:", ids[:5])

	storyCh, done := storyFetcherPool(c.client, ids, c.workers)
	stories := retrieveStories(c.stories, storyCh, done)
	c.mutex.Lock()
	c.items = stories
	c.lastUpdated = time.Now()
	c.mutex.Unlock()

	for i := 0; i < 5; i++ {
		fmt.Println("id:", c.items[i].ID)
		fmt.Println("title:", c.items[i].Title)
	}

	return nil
}

func (c *cache) UpdateEvery(d time.Duration) {
	tick := time.NewTicker(d)
	defer tick.Stop()
	go func() {
		for {
			<-tick.C
			err := c.Update()
			if err != nil {
				log.Fatalln("Couldn't update cache :(")
			}
		}
	}()
}

func (c *cache) Get() []item {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.items
}
