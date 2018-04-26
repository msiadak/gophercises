package main

import (
	"fmt"
	"strings"

	"github.com/msiadak/gophercises/deck"
)

const (
	minPlayers = 2
	maxPlayers = 8
)

type hand struct {
	cards []deck.Card
}

func (h *hand) Draw(d *[]deck.Card, n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("Cards before: %s\n", h)
		h.cards = append(h.cards, (*d)[0])
		fmt.Printf("Cards after: %s\n", h)
		*d = (*d)[1:]
	}
}

func (h hand) String() string {
	var cardStrings []string
	for _, card := range h.cards {
		cardStrings = append(cardStrings, card.String())
	}
	return strings.Join(cardStrings, ", ")
}

func (h hand) DealerString() string {
	cardStrings := []string{"*HIDDEN*"}
	for _, card := range h.cards[1:] {
		cardStrings = append(cardStrings, card.String())
	}
	return strings.Join(cardStrings, ", ")
}

func deal(gameDeck *[]deck.Card, nPlayers int) []hand {
	if nPlayers < minPlayers || nPlayers > maxPlayers {
		panic("Wrong number of players")
	}
	hands := make([]hand, nPlayers)
	for i := 0; i < 2; i++ {
		for _, hand := range hands {
			fmt.Printf("Cards before: %s\n", hand)
			hand.Draw(gameDeck, 1)
			fmt.Printf("Cards after: %s\n", hand)
		}
	}
	return hands
}

func main() {
	gameDeck := deck.New(deck.Shuffle)

	hands := deal(&gameDeck, 8)

	for i, hand := range hands[:len(hands)-1] {
		fmt.Printf("Player %d: %s\n", i+1, hand)
	}
	fmt.Printf("Dealer: %s\n", hands[len(hands)-1])
}
