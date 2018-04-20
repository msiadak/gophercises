package link

import (
	"os"
	"strings"
	"testing"
)

func equalLinks(a, b []Link) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestExtractLinks(t *testing.T) {
	t.Run("Return an empty slice for an empty input", func(t *testing.T) {
		want := make([]Link, 0)
		got, err := ExtractLinks(strings.NewReader(""))
		if err != nil {
			t.Error(err)
		}

		if !equalLinks(want, got) {
			t.Errorf("want: '%v', got: '%v'", want, got)
		}
	})

	t.Run("Returns one link", func(t *testing.T) {
		want := []Link{Link{"http://web.archive.org", "Wayback Machine"}}
		got, err := ExtractLinks(strings.NewReader(`<a href="http://web.archive.org">Wayback Machine</a>`))
		if err != nil {
			t.Fatal(err)
		}

		if !equalLinks(want, got) {
			t.Errorf("want: '%v', got: '%v'", want, got)
		}
	})

	t.Run("Returns two links", func(t *testing.T) {
		want := []Link{
			Link{"http://a.com", "a"},
			Link{"http://b.com", "b"},
		}
		got, err := ExtractLinks(strings.NewReader(`
			<a href="http://a.com">a</a>
			<a href="http://b.com">b</a>
		`))
		if err != nil {
			t.Fatal(err)
		}

		if !equalLinks(want, got) {
			t.Errorf("want: '%v', got: '%v'", want, got)
		}
	})

	t.Run("Parses ex1 correctly", func(t *testing.T) {
		want := []Link{
			Link{"/other-page", "A link to another page"},
		}
		file, err := os.Open("examples/ex1.html")
		if err != nil {
			t.Fatal(err)
		}

		got, err := ExtractLinks(file)
		if err != nil {
			t.Fatal(err)
		}

		if !equalLinks(want, got) {
			t.Errorf("want: '%v', got: '%v'", want, got)
		}
	})

	t.Run("Parses ex2 correctly", func(t *testing.T) {
		want := []Link{
			Link{
				"https://www.twitter.com/joncalhoun",
				"Check me out on twitter",
			},
			Link{
				"https://github.com/gophercises",
				"Gophercises is on Github !",
			},
		}
		file, err := os.Open("examples/ex2.html")
		if err != nil {
			t.Fatal(err)
		}

		got, err := ExtractLinks(file)
		if err != nil {
			t.Fatal(err)
		}

		if !equalLinks(want, got) {
			t.Errorf("want: '%v', got: '%v'", want, got)
		}
	})

	t.Run("Parses ex3 correctly", func(t *testing.T) {
		want := []Link{
			Link{"#", "Login"},
			Link{"/lost", "Lost? Need help?"},
			Link{"https://twitter.com/marcusolsson", "@marcusolsson"},
		}
		file, err := os.Open("examples/ex3.html")
		if err != nil {
			t.Fatal(err)
		}

		got, err := ExtractLinks(file)
		if err != nil {
			t.Fatal(err)
		}

		if !equalLinks(want, got) {
			t.Errorf("want: '%v', got: '%v'", want, got)
		}
	})

	t.Run("Parses ex4 correctly", func(t *testing.T) {
		want := []Link{
			Link{"/dog-cat", "dog cat"},
		}
		file, err := os.Open("examples/ex4.html")
		if err != nil {
			t.Fatal(err)
		}

		got, err := ExtractLinks(file)
		if err != nil {
			t.Fatal(err)
		}

		if !equalLinks(want, got) {
			t.Errorf("want: '%v', got: '%v'", want, got)
		}
	})
}
