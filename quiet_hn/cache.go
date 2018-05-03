package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/msiadak/gophercises/quiet_hn/hn"
)

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
	ids, err := c.client.TopItems()
	if err != nil {
		return fmt.Errorf("Failed to load top stories: %s", err)
	}

	storyCh, done := storyFetcherPool(c.client, ids, c.workers)
	stories := retrieveStories(c.stories, storyCh, done)
	c.mutex.Lock()
	c.items = stories
	sort.Sort(byIDSlice{ids, c.items})
	c.lastUpdated = time.Now()
	c.mutex.Unlock()

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
