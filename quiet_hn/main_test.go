package main

import (
	"html/template"
	"net/http/httptest"
	"testing"

	"github.com/msiadak/gophercises/link"
	"github.com/msiadak/gophercises/quiet_hn/hn"
)

func TestHandler(t *testing.T) {
	t.Run("Returns 32 links", func(t *testing.T) {
		c := newCache(hn.Client{}, 16, 30)
		c.Update()

		tpl := template.Must(template.ParseFiles("./index.gohtml"))

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()
		handler(c, tpl)(rr, req)

		resp := rr.Result()
		links, err := link.ExtractLinks(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if len(links) != 32 {
			t.Fatalf("Want: 32 links; Got: %d links", len(links))
		}
	})
}
