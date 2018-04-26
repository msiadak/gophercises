package deck

import (
	"fmt"
	"math/rand"
	"sort"
)

// Possible suits in a standard 52-card deck.
const (
	Spades Suit = iota + 1
	Diamonds
	Clubs
	Hearts
)

// Possible ranks in a standard 52-card deck.
const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Joker Rank = -1
)

// Suit represents a card's suit.
type Suit int

//go:generate stringer -type=Suit

// Rank represents a card's numeric or face value.
type Rank int

//go:generate stringer -type=Rank

// Card represents a playing card in a standard 52-card deck.
type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	var suffix string
	if c.Suit > 0 {
		suffix = fmt.Sprintf(" of %s", c.Suit)
	}

	return fmt.Sprintf("%s%s", c.Rank, suffix)
}

// New returns a new deck (slice of cards).
//
// By default, the deck is returned in the order that a new deck of cards will
// come in the package (in ascending value by suit, then by rank as shown
// in the order of the Suit and Rank constants above).
func New(options ...option) []Card {
	cards := make([]Card, 0, 52)

	for i := 1; i <= 4; i++ {
		for j := 1; j <= 13; j++ {
			cards = append(cards, Card{Suit(i), Rank(j)})
		}
	}

	for _, option := range options {
		cards = option(cards)
	}

	return cards
}

type option func([]Card) []Card

type cardSorter struct {
	cards []Card
	by    func(c1, c2 *Card) bool
}

func (s *cardSorter) Len() int {
	return len(s.cards)
}

func (s *cardSorter) Swap(i, j int) {
	s.cards[i], s.cards[j] = s.cards[j], s.cards[i]
}

func (s *cardSorter) Less(i, j int) bool {
	return s.by(&s.cards[i], &s.cards[j])
}

// SortBy accepts a "by" function (should return true if c1 should be sorted
// before c2) and returns a closure that operates on a slice of cards.
//
// The returned closure will sort the deck of cards based on the provided "by"
// function and can be passed to New() as an option to sort the new deck.
func SortBy(by func(c1, c2 *Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sorted := make([]Card, len(cards))
		copy(sorted, cards)
		sort.Sort(&cardSorter{sorted, by})
		return sorted
	}
}

// SortReverse returns a function that sorts a slice of cards in descending
// order of Suit, then Rank.
func SortReverse() func([]Card) []Card {
	return SortBy(func(c1, c2 *Card) bool {
		return c1.Suit > c2.Suit || (c1.Suit == c2.Suit && c1.Rank > c2.Rank)
	})
}

// Shuffle randomizes the order of a slice of cards.
func Shuffle(cards []Card) []Card {
	shuffled := make([]Card, len(cards))
	copy(shuffled, cards)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

// Jokers returns a function that adds n jokers to a slice of cards.
func Jokers(n uint) func([]Card) []Card {
	return func(cards []Card) []Card {
		newCards := make([]Card, len(cards), len(cards)+int(n))
		copy(newCards, cards)
		for i := uint(0); i < n; i++ {
			newCards = append(newCards, Card{Rank: Joker})
		}
		return newCards
	}
}

// Filter returns a closure that filters cards from a card slice.
func Filter(filterCards []Card) func([]Card) []Card {
	return func(cards []Card) []Card {
		filtered := make([]Card, 0)
		for _, c := range cards {
			for _, f := range filterCards {
				switch {
				case c.Suit == f.Suit && f.Rank == 0:
				case c.Rank == f.Rank && f.Suit == 0:
				case c.Suit == f.Suit && c.Rank == f.Rank:
				default:
					filtered = append(filtered, c)
				}
			}
		}
		return filtered
	}
}
