package main

import (
	"github.com/gdamore/tcell"
	// "log"
	// "fmt"
)

///////////////////////////////////////////////////////////////////////////////
// Data
///////////////////////////////////////////////////////////////////////////////

// CardValue is an enum denoting the different possible
// values of a playing card.
type CardValue int

const (
	NoneValue CardValue = iota
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

// This function is a workaround to get a constant global array
func getValueToString() []string {
	return []string{"None", "Ace", "Two", "Three", "Four", "Five",
		"Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
}

// NUM_VALUES is the number of different possible values
// (not including NoneValue).
var NUM_VALUES = int(King)

// CardSuit is an enum denoting the possible suits of a playing card.
type CardSuit int

const (
	NoneSuit CardSuit = iota
	Spades
	Hearts
	Clubs
	Diamonds
)

// This function is a workaround to get a global constant array
func getSuitToString() []string {
	return []string{"None", "Spades", "Hearts", "Clubs", "Diamonds"}
}

// Card represents a single playing card.
type Card struct {
	suit  CardSuit
	value CardValue
}

///////////////////////////////////////////////////////////////////////////////
// Card Functions
///////////////////////////////////////////////////////////////////////////////

func (card Card) toString() string {
	return getValueToString()[card.value] +
		" of " + getSuitToString()[card.suit]
}
func (cv CardValue) toString() string {
	return getValueToString()[cv]
}

// isBlank returns true iff the card has the NoneSuit and NoneValue.
func (card Card) isBlank() bool {
	if card.suit == NoneSuit || card.value == NoneValue {
		return true
	}
	return false
}

// Render renders a single face-up card on the terminal screen, with the
// upper-left corner at point x, y.
func (card Card) Render(s tcell.Screen, x int, y int) {
	color := tcell.ColorGreen

	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(color)
	box1 := Box{s, x, y, x + CARD_WIDTH, y + CARD_HEIGHT, boxStyle,
		card.toString(), false}
	box1.Draw()
}

// RenderFlipped renders a face-down card on the terminal screen, with
// the upper-left corner at point x, y.
func (card Card) RenderFlipped(s tcell.Screen, x int, y int) {
	color := tcell.ColorGreen

	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(color)
	box1 := Box{s, x, y, x + CARD_WIDTH, y + CARD_HEIGHT,
		boxStyle, "", false}
	box1.Draw()
}
