package deck

import (
	"testing"
)

var suitOrder string
var rankOrder string

func TestNew(t *testing.T) {
	t.Run("Default deck", func(t *testing.T) {
		got := New()
		for i := 0; i < len(got)-1; i++ {
			c1, c2 := got[i], got[i+1]
			if (c1.Suit > c2.Suit) ||
				(c1.Suit == c2.Suit && c1.Rank > c2.Rank) {
				t.Errorf("cards[%d]: %s; cards[%d]: %s", i, c1, i+1, c2)
				t.Fatal("expected suit and rank to be sorted in ascending order")
			}
		}
	})
	t.Run("Descending order", func(t *testing.T) {
		desc := func(c1, c2 *Card) bool {
			return c1.Suit > c2.Suit || (c1.Suit == c2.Suit && c1.Rank > c2.Rank)
		}
		got := New(SortBy(desc))
		for i := 0; i < len(got)-1; i++ {
			c1, c2 := got[i], got[i+1]
			if (c1.Suit < c2.Suit) ||
				(c1.Suit == c2.Suit && c1.Rank < c2.Rank) {
				t.Errorf("cards[%d]: %v; cards[%d]: %v", i, c1, i+1, c2)
				t.Fatal("expected suit and rank to be sorted in descending order")
			}
		}
	})
}
