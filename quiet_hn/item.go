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

type fetchItemJob struct {
	ID         int
	ResponseCh chan item
}

func fetchStories(client hn.Client, ids []int, numWorkers, numStories int) []item {
	jobs := make(chan fetchItemJob, numWorkers)
	responses := make(chan chan item, numWorkers)
	done := make(chan struct{})
	defer close(done)
	defer close(jobs)
	defer close(responses)

	// Spawn workers to handle fetch jobs
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				case job := <-jobs:
					hnItem, err := client.GetItem(job.ID)
					if err != nil {
						log.Printf("Couldn't get item %d: %s\n", job.ID, err)
						continue
					}
					item := parseHNItem(hnItem)
					job.ResponseCh <- item
				}
			}
		}()
	}

	// Feed the jobs to the jobs chan and pass on the response chans
	go func() {
		for _, id := range ids {
			select {
			case <-done:
				return
			default:
				job := fetchItemJob{id, make(chan item)}
				responses <- job.ResponseCh
				jobs <- job
			}
		}
	}()

	// Receive the stories and store them in a slice
	stories := make([]item, numStories)
	for i := 0; i < numStories; i++ {
		for {
			responseCh := <-responses
			item := <-responseCh
			close(responseCh)
			if isStoryLink(item) {
				stories[i] = item
				break
			}
		}
	}

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
