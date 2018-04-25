package deck

import "sort"

// Possible suits in a standard 52-card deck.
const (
	Spades Suit = iota
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
)

// Suit represents a card's suit.
type Suit int

// Rank represents a card's numeric or face value.
type Rank int

// Card represents a playing card in a standard 52-card deck.
type Card struct {
	Suit
	Rank
}

// New returns a new deck (slice of cards).
//
// By default, the deck is returned in the order that a new deck of cards will
// come in the package (in ascending value by suit, then by rank as shown
// in the order of the Suit and Rank constants above).
func New(options ...func([]Card)) []Card {
	cards := make([]Card, 52)

	for i := 0; i < 4; i++ {
		for j := 1; j <= 13; j++ {
			cards = append(cards, Card{Suit(i), Rank(j)})
		}
	}

	for _, option := range options {
		option(cards)
	}

	return cards
}

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
func SortBy(by func(c1, c2 *Card) bool) func([]Card) {
	return func(cards []Card) {
		sort.Sort(&cardSorter{cards, by})
	}
}

