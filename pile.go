package main

import (
	"github.com/gdamore/tcell"
	// "fmt"
)

// Pile is one of the ten stacks of partially visible cards
type Pile struct {
	visible Deck;
	invisible Deck;
}

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

func (pile Pile) Height() int {
	return pile.invisible.Size() + pile.visible.Size()*2
}

// if there are not n visible cards, return the None card
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


