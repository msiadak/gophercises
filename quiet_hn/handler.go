package main

import (
	"html/template"
	"net/http"
	"time"
)

type templateData struct {
	Stories []item
	Time    time.Duration
}

func handler(c *cache, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		data := templateData{
			Stories: c.Get(),
			Time:    time.Now().Sub(start),
		}
		err := tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}
