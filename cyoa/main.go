package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Arc struct {
	Title   string
	Story   []string
	Options []Option
}

type Option struct {
	Text string
	Arc  string
}

type AdventureHandler struct {
	Arcs     map[string]Arc
	Template template.Template
}

func (a AdventureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	fmt.Printf("Path requested: %s\n", path)
	if arc, ok := a.Arcs[path]; ok {
		a.Template.Execute(w, arc)
		return
	}
	http.Redirect(w, r, "/intro", 302)
}

func parseJSON(filename string) (map[string]Arc, error) {
	f, err := os.Open("gopher.json")
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(f)
	arcs := make(map[string]Arc)
	err = decoder.Decode(&arcs)
	if err != nil {
		return nil, err
	}

	return arcs, nil
}

func main() {
	arcs, err := parseJSON("gopher.json")
	if err != nil {
		panic(err)
	}

	template, err := template.ParseFiles("arc.html")
	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":8080", AdventureHandler{arcs, *template}))
}
