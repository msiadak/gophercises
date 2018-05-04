package main

import (
	"fmt"
	"log"
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
	log.Println("Starting cache update")
	ids, err := c.client.TopItems()
	if err != nil {
		return fmt.Errorf("Failed to load top stories: %s", err)
	}

	stories, err := fetchStories(c.client, ids, c.workers, c.stories)
	if err != nil {
		return err
	}

	c.mutex.Lock()
	c.items = stories
	c.lastUpdated = time.Now()
	c.mutex.Unlock()

	log.Println("Finished cache update")
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
