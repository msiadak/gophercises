package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
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
