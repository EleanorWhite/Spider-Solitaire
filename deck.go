package main

import (
	"github.com/gdamore/tcell"
	"math/rand"
	"time"
	"log"
	// "fmt"
)


// Deck represents a deck of playing cards, where cards
// are drawn starting from the highest index.
// The "top" of the deck is defined as the next ones to be drawn.
type Deck struct {
	cards []Card
}



///////////////////////////////////////////////////////////////////////////////
// Deck Functions
///////////////////////////////////////////////////////////////////////////////

// IsEmpty return true iff deck has no cards.
func (deck Deck) IsEmpty() bool {
	return len(deck.cards) == 0
}

// Size returns the number of cards in deck.
func (deck Deck) Size() int {
	return len(deck.cards);
}

// Add puts card on the top of deck.
func (deck *Deck) Add(card Card) {
	if (card.isBlank()) { return }
	deck.cards = append(deck.cards, card);
}

// Draw takes a card off the top of deck and returns it.
func (deck *Deck) Draw() Card{
	if (len(deck.cards) == 0) {
		return Card{NoneSuit, NoneValue}
	}
	card := deck.cards[len(deck.cards)-1]
	deck.cards = deck.cards[:len(deck.cards)-1]
	return card
}

// PeekTopNCards returns a slice of the top n cards of Deck deck..
// The Cards are not removed from deck.
func (deck Deck) PeekTopNCards(n int) []Card{
	length := len(deck.cards)
	if (length < n) {
		log.Fatalf("Asking for more cards than exist")
	}
	return deck.cards[length - n:]
}

// GetTopNCards returns a slice of a new array containing
// a copy of the top n cards of Deck deck.
// The Cards are removed from deck.
func (deck *Deck) GetTopNCards(n int) []Card{
	length := len(deck.cards)
	if (length < n) {
		log.Fatalf("Asking for more cards than exist")
	}
	topN := deck.cards[length - n:]
	newTopN := make([]Card, len(topN))
	copy(newTopN, topN)

	deck.cards = deck.cards[:length - n]
	return newTopN
}

// PeekNthCard returns the nth Card from the top of deck.
// The Card is not removed from deck.
func (deck Deck) PeekNthCard(n int) Card {
	// fmt.Println("n: ", n)
	if (len(deck.cards) <= n) {
		return Card{NoneSuit, NoneValue}
	}
	return deck.cards[len(deck.cards)- 1 - n]
}

// NewDeck creates and returns a new empty Deck. 
// initSize is the initial size of the underlying array
// for the Deck, but the Deck will be resized to accommodate
// cards as needed.
func NewDeck(initSize int) Deck {
	var deck Deck
	deck.cards = make([]Card, 0, initSize)
	return deck
}

// Swap swaps the cards in the indices fst and snd or the cards array.
func (deck *Deck) Swap(fst int, snd int) {
	deck.cards[fst], deck.cards[snd] = deck.cards[snd], deck.cards[fst]
}

// Shuffle randomizes the order of deck.
func (deck *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck.cards), deck.Swap)
}


///////////////////////////////////////////////////////////////////////////////
// Debug Utilities
///////////////////////////////////////////////////////////////////////////////

// ToString returns a string representing deck.
func (deck *Deck) ToString() string {
	if (len(deck.cards) == 0) { return "" }
	var str string = ""
	for _, v := range deck.cards {
		str += v.toString() + ", "
	}
	str = str[:len(str) -1]
	return str
}


///////////////////////////////////////////////////////////////////////////////
// Graphics Function
///////////////////////////////////////////////////////////////////////////////

// Render renders a deck of cards face side up
func (deck Deck) Render(s tcell.Screen, x int, y int, selected bool) {
	for i, v := range deck.cards {
		v.Render(s, x, y + 2*i, selected)
	}
}
// RenderFlipped renders a deck of cards face side down
func (deck Deck) RenderFlipped(s tcell.Screen, x int, y int, selected bool) {
	for i, v := range deck.cards {
		v.RenderFlipped(s, x, y+i, selected)
	}
}
