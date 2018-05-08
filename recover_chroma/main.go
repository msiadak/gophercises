package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", srcMw(devMw(mux))))
}

func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, string(stack))
			}
		}()
		app.ServeHTTP(w, r)
	}
}

func srcMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !(strings.HasPrefix(r.URL.Path, "/src/") && strings.HasSuffix(r.URL.Path, ".go")) {
			app.ServeHTTP(w, r)
			return
		}
		path := r.URL.Path[len("/src/"):]

		buf, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			http.Error(w, fmt.Sprintf("Couldn't open file: %s\n", path), http.StatusNotFound)
			return
		}

		text := html.EscapeString(string(buf))
		fmt.Fprintf(w, "<h1>%s</h1><pre>%s</pre>", path, text)
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
