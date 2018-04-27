package deck

import (
	"fmt"
	"testing"
)

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
		got := New(SortReverse())
		for i := 0; i < len(got)-1; i++ {
			c1, c2 := got[i], got[i+1]
			if (c1.Suit < c2.Suit) ||
				(c1.Suit == c2.Suit && c1.Rank < c2.Rank) {
				t.Errorf("cards[%d]: %v; cards[%d]: %v", i, c1, i+1, c2)
				t.Fatal("expected suit and rank to be sorted in descending order")
			}
		}
	})

	t.Run("Shuffle", func(t *testing.T) {
		got := New(Shuffle)[:5]
		want := []Card{
			Card{Clubs, King},
			Card{Clubs, Ace},
			Card{Hearts, Five},
			Card{Hearts, King},
			Card{Diamonds, Ace},
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("at index %d; got: %s; want: %s", i, got[i], want[i])
			}
		}
	})

	t.Run("Jokers", func(t *testing.T) {
		var want uint = 5
		var got uint
		cards := New(Jokers(want))

		for _, c := range cards {
			if c.Rank == Joker {
				got++
			}
		}

		if got < want {
			t.Errorf("wanted %d Jokers, got %d", want, got)
		}
	})

	t.Run("Filter", func(t *testing.T) {
		t.Run("by Rank", func(t *testing.T) {
			filter := []Card{Card{Rank: Ace}}
			cards := New(Filter(filter))

			for _, c := range cards {
				if c.Rank == filter[0].Rank {
					t.Errorf("expected no %s cards; got %s", filter[0].Rank, c)
				}
			}
		})

		t.Run("by Suit", func(t *testing.T) {
			filter := []Card{Card{Suit: Clubs}}
			cards := New(Filter(filter))

			for _, c := range cards {
				if c.Suit == filter[0].Suit {
					t.Errorf("expected no %s cards; got %s", filter[0].Suit, c)
				}
			}
		})

		t.Run("by Rank and Suit", func(t *testing.T) {
			filter := []Card{Card{Hearts, King}}
			cards := New(Filter(filter))

			for _, c := range cards {
				if c == filter[0] {
					t.Errorf("expected no %s; got %s", c, filter[0])
				}
			}
		})
	})

	t.Run("Multiply", func(t *testing.T) {
		cards := New(Multiply(3))
		got := len(cards)
		want := 156

		if want != got {
			t.Errorf("Wanted %d cards, got %d", want, got)
		}
	})
}

func TestCard(t *testing.T) {
	cards := []Card{
		Card{Spades, Ace},
		Card{Diamonds, King},
		Card{Clubs, Queen},
		Card{Hearts, Jack},
		Card{Rank: Joker},
	}
	for _, card := range cards {
		fmt.Println(card)
	}

}
