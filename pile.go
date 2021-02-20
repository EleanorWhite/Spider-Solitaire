package main

import (
	"github.com/gdamore/tcell"
	// "fmt"
)

///////////////////////////////////////////////////////////////////////////////
// Data Types
///////////////////////////////////////////////////////////////////////////////


// Pile is one of the ten stacks of partially visible cards
type Pile struct {
	visible Deck;
	invisible Deck;
}


///////////////////////////////////////////////////////////////////////////////
// Functions to modify and get info about Piles
///////////////////////////////////////////////////////////////////////////////


// TopNMovable returns true if the top n cards in the visible
// part of the Pile can be moved together. Cards can be moved
// together if the cards are in ascending order starting from 
// the top.
func (pile Pile) TopNMovable(n int) bool {
	// fmt.Println(pile.visible.ToString())
	if (n > pile.visible.Size()) {
		return false
	}
	var cards []Card = pile.visible.PeekTopNCards(n)
	for i, v := range cards {
		if (i+1 < len(cards)) {
			//fmt.Println(v.value, ", ", cards[i+1].value)
			if (v.value - 1 != cards[i+1].value) {
				return false
			}
		}
	}
	return true
}

// PeekNthCard returns the nth card from the top of the 
// visible part of the pile. It returns a Card with NoneValue
// and NoneSuit if there is not an nth card.
func (pile Pile) PeekNthCard(n int) Card {
	return pile.visible.PeekNthCard(n)
}

// GetTopNCards returns a slice of the top n cards of the 
// visible portion of a Pile.
// Top is defined by the next ones to be drawn.
// The Cards are removed.
func (pile *Pile) GetTopNCards(n int) []Card{
	moved := pile.visible.GetTopNCards(n)
	if (pile.visible.IsEmpty()) {
		pile.visible.Add(pile.invisible.Draw())
	}
	return moved
}

// IsEmpty returns true iff there are no cards in the pile 
// (visible or invisible).
func (pile Pile) IsEmpty() bool {
	return pile.visible.IsEmpty() && pile.invisible.IsEmpty()
}

///////////////////////////////////////////////////////////////////////////////
// Graphics
///////////////////////////////////////////////////////////////////////////////

// Render renders a pile on the terminal
func (pile Pile) Render(s tcell.Screen, x int, y int, selected bool) {
	pile.invisible.RenderFlipped(s, x, y, selected)
	pile.visible.Render(s, x, y + pile.invisible.Size(), selected)
}

// Height returns the height of the pile on the screen. 
// The unit is the y coordinate used by the Box struct.
func (pile Pile) Height() int {
	if (pile.IsEmpty()) {
		// empty pile highlights should line up with the bottom
		// card of the pile.
		return 2  
	}
	return pile.invisible.Size() + pile.visible.Size()*2
}


