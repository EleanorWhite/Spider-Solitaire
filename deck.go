package main

import (
	"github.com/gdamore/tcell"
	"math/rand"
	"time"
	"log"
	// "fmt"
)

type CardValue int

const  (
	NoneValue CardValue   = iota
	Ace
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
// TODO: MAKE CONST
var ValueToString = []string{"None", "Ace", "Two", "Three", "Four", "Five", 
	"Six", "Seven", "Eight", "Nine", "Ten", "Jack","Queen","King"}
var NUM_VALUES = len(ValueToString) - 1


type CardSuit int

const (
	NoneSuit  CardSuit = iota
	Spades  
	Hearts  
	Clubs   
	Diamonds
)

var SuitToString = []string{"None", "Spades", "Hearts", "Clubs", "Diamonds"}

type Card struct {
	suit CardSuit
	value CardValue
}

type Deck struct {
	cards []Card
}

///////////////////////////////////////////////////////////////////////////////
// Card Stuff
///////////////////////////////////////////////////////////////////////////////

func (c Card) toString() string {
	return ValueToString[c.value] + " of " + SuitToString[c.suit]
}
func (cv CardValue) toString() string {
	return ValueToString[cv] 
}

func (card Card) isBlank() bool {
	if(card.suit == NoneSuit || card.value == NoneValue) { 
		return true
	}
	return false
}

func (card Card) Render(s tcell.Screen, x int, y int, selected bool) {
	color := tcell.ColorGreen;
	if selected {
		color = tcell.ColorYellow
	}
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(color)
	box1 := Box{s, x, y, x+CARD_WIDTH, y+CARD_HEIGHT, boxStyle, 
		card.toString(), false};
	box1.Draw();
}


func (card Card) RenderFlipped(s tcell.Screen, x int, y int, selected bool) {
	color := tcell.ColorGreen;
	if selected {
		color = tcell.ColorYellow
	}
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(color)
	box1 := Box{s, x, y, x+CARD_WIDTH, y+CARD_HEIGHT, 
		boxStyle, "", false};
	box1.Draw();
}

///////////////////////////////////////////////////////////////////////////////
// Deck Stuff
///////////////////////////////////////////////////////////////////////////////

func (deck Deck) IsEmpty() bool {
	return len(deck.cards) == 0
}

func (deck Deck) Size() int {
	return len(deck.cards);
}

func (deck *Deck) Add(card Card) {
	if (card.isBlank()) { return }
	deck.cards = append(deck.cards, card);
}

func (deck *Deck) Draw() Card{
	if (len(deck.cards) == 0) {
		return Card{NoneSuit, NoneValue}
	}
	card := deck.cards[len(deck.cards)-1]
	deck.cards = deck.cards[:len(deck.cards)-1]
	return card
}

// func (deck *Deck) GetCardAtIndex(index int) Card{
// 	card := deck.cards[index]
// 	deck.cards[index] = deck.cards[len(deck.cards)-1]
// 	deck.cards = deck.cards[:len(deck.cards)-1]
// 	return card
// }

// PeekTopNCards returns a slice of the top n cards of a Deck.
// Top is defined by the next ones to be drawn.
// The Cards are not removed.
func (deck *Deck) PeekTopNCards(n int) []Card{
	length := len(deck.cards)
	if (length < n) {
		log.Fatalf("Asking for more cards than exist")
	}
	return deck.cards[length - n:]
}

// GetTopNCards returns a slice of a new array containing
// a copy of the top n cards of a Deck.
// Top is defined by the next ones to be drawn.
// The Cards are removed.
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

func (deck *Deck) PeekNthCard(n int) Card {
	// fmt.Println("n: ", n)
	if (len(deck.cards) <= n) {
		return Card{NoneSuit, NoneValue}
	}
	return deck.cards[len(deck.cards)- 1 - n]
}

func NewDeck(maxSize int) Deck {
	var deck Deck
	deck.cards = make([]Card, 0, maxSize)
	return deck
}

func (deck *Deck) ToString() string {
	if (len(deck.cards) == 0) { return "" }
	var str string = ""
	for _, v := range deck.cards {
		str += v.toString() + ", "
	}
	str = str[:len(str) -1]
	return str
}

func (deck *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck.cards), deck.Swap)
}

func (deck *Deck) Swap(fst int, snd int) {
	deck.cards[fst], deck.cards[snd] = deck.cards[snd], deck.cards[fst]
}

func (deck Deck)CardInDeck(card Card) bool {
    for _, v := range deck.cards {
        if v == card {
            return true
        }
    }
    return false
}


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
