package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/msiadak/gophercises/link"
)

func TestHandler(t *testing.T) {
	t.Run("Returns 32 links", func(t *testing.T) {
		tpl := template.Must(template.ParseFiles("./index.gohtml"))
		handler := handler(30, 8, tpl)

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)

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

func benchmarkHandler(numStories int, numWorkers int, b *testing.B) {
	tpl := template.Must(template.ParseFiles("./index.gohtml"))
	server := httptest.NewServer(handler(numStories, numWorkers, tpl))

	resp, err := http.Get(server.URL)
	if err != nil {
		b.Fatal(err)
	}

	if resp.StatusCode != 200 {
		b.Fatal(resp.Status)
	}
}

func BenchmarkHandler30_1(b *testing.B)  { benchmarkHandler(30, 1, b) }
func BenchmarkHandler30_2(b *testing.B)  { benchmarkHandler(30, 2, b) }
func BenchmarkHandler30_4(b *testing.B)  { benchmarkHandler(30, 4, b) }
func BenchmarkHandler30_8(b *testing.B)  { benchmarkHandler(30, 8, b) }
func BenchmarkHandler30_16(b *testing.B) { benchmarkHandler(30, 16, b) }
func BenchmarkHandler30_30(b *testing.B) { benchmarkHandler(30, 32, b) }
