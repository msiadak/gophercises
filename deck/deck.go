package deck

type Suit int

const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts
)

type Value int

const (
	Ace Value = iota + 1
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

type Card struct {
	Suit  Suit
	Value Value
}

func New(options ...func([]Card)) []Card {
	cards := make([]Card, 52)

	for i := 0; i < 4; i++ {
		for j := 1; j <= 13; j++ {
			cards = append(cards, Card{Suit(i), Value(j)})
		}
	}

	for _, option := range options {
		option(cards)
	}

	return cards
}

func Sort([]Card