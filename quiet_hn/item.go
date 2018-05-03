package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/msiadak/gophercises/quiet_hn/hn"
)

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
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

type byIDSlice struct {
	ids   []int
	items []item
}

func (s byIDSlice) Len() int {
	return len(s.items)
}

func (s byIDSlice) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s byIDSlice) Less(i, j int) bool {
	return s.findIDIndex(s.items[i].ID) < s.findIDIndex(s.items[j].ID)
}

func (s byIDSlice) findIDIndex(id int) int {
	for i, v := range s.ids {
		if v == id {
			return i
		}
	}
	return -1
}
