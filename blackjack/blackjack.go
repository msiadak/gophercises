package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/msiadak/gophercises/deck"
)

const (
	minPlayers = 2
	maxPlayers = 8
)

type hand struct {
	Cards []deck.Card
}

type handValue struct {
	value int
	soft  bool
}

func (h *hand) Draw(d *[]deck.Card, n int) {
	for i := 0; i < n; i++ {
		h.Cards = append(h.Cards, (*d)[0])
		*d = (*d)[1:]
	}
}

func (h hand) String() string {
	var cardStrings []string
	for _, card := range h.Cards {
		cardStrings = append(cardStrings, card.String())
	}
	return strings.Join(cardStrings, ", ")
}

func (h hand) DealerString() string {
	cardStrings := []string{"??"}
	for _, card := range h.Cards[1:] {
		cardStrings = append(cardStrings, card.String())
	}
	return strings.Join(cardStrings, ", ")
}

func (h hand) Value() handValue {
	hv := handValue{}
	acesLast := deck.SortBy(func(c1, c2 *deck.Card) bool {
		return c1.Rank < c2.Rank && c1.Rank != deck.Ace
	})
	for _, c := range acesLast(h.Cards) {
		switch c.Rank {
		case deck.Ace:
			if hv.value+11 > 21 {
				hv.value++
			} else {
				hv.value += 11
				hv.soft = true
			}
		case deck.Ten, deck.Jack, deck.Queen, deck.King:
			hv.value += 10
		default:
			hv.value += int(c.Rank)
		}
	}
	return hv
}

func (h hand) Busted() bool {
	return h.Value().value > 21
}

func (h hand) LastCard() deck.Card {
	return h.Cards[len(h.Cards)-1]
}

func deal(gameDeck *[]deck.Card, nPlayers int) []hand {
	if nPlayers < minPlayers || nPlayers > maxPlayers {
		panic("Wrong number of players")
	}
	hands := make([]hand, nPlayers)
	for i := 0; i < 2; i++ {
		for j := range hands {
			hands[j].Draw(gameDeck, 1)
		}
	}
	return hands
}

func main() {
	rand.Seed(time.Now().Unix())
	gameDeck := deck.New(deck.Shuffle)

	hands := deal(&gameDeck, 2)
	dealerHand := hands[len(hands)-1]
	playerHands := hands[:len(hands)-1]

	for i, hand := range playerHands {
		fmt.Printf("Player %d: %s\n", i+1, hand)
	}
	fmt.Printf("Dealer: %s\n\n", dealerHand.DealerString())

	for i := range playerHands {
		fmt.Printf("Player %d's Turn:\n", i+1)
		var act string
	actionLoop:
		for {
			fmt.Printf("Player Cards: %s\n", playerHands[i])
			fmt.Printf("Dealer Cards: %s\n", dealerHand.DealerString())
			fmt.Print("(H)it or (S)tand? ")
			fmt.Scanf("%s", &act)
			switch strings.ToLower(act) {
			case "h":
				playerHands[i].Draw(&gameDeck, 1)
				fmt.Printf("Hit: %s\n", playerHands[i].LastCard())
				if playerHands[i].Busted() {
					fmt.Printf("Busted! (Score: %d)\n", playerHands[i].Value().value)
					break actionLoop
				}
			case "s":
				fmt.Println("You stood")
				fmt.Printf("Score: %d\n", playerHands[i].Value().value)
				break actionLoop
			default:
				fmt.Println("Invalid input (try H or S)")
			}
		}
	}

	fmt.Printf("Dealer's Hand: %s\n", dealerHand)
	for dealerHand.Value().value <= 16 ||
		(dealerHand.Value().value == 17 && dealerHand.Value().soft) {
		dealerHand.Draw(&gameDeck, 1)
		fmt.Printf("Dealer hits: %s\n", dealerHand.LastCard())
	}
	if dealerHand.Busted() {
		fmt.Printf("Dealer busts! (Value: %d)\n", dealerHand.Value().value)
	} else {
		fmt.Printf("Dealer has %d\n", dealerHand.Value().value)
	}

}
