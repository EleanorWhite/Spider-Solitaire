package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell"
)

const NUM_PILES = 10
const CARD_WIDTH = 11
const CARD_HEIGHT = 7

///////////////////////////////////////////////////////////////////////////////
// Data Types
///////////////////////////////////////////////////////////////////////////////

// Game contains all info about the current game
type Game struct {
	deck        Deck // The remaining deck which has cards not yet on piles
	piles       [NUM_PILES]Pile
	highlighted Selected // which card the cursor is over
	toMove      bool     // whether the user has cards selected that they might move
	selected    Selected // which card(s) are selected
}

// Selected is a description of cards currently selected/highlighted
// by the user
type Selected struct {
	x        int // which pile is highlighted
	y        int // whether deck is highlighted (0), or pile is highlighted (1)
	numCards int // How many cards in a pile are highlighted
}

///////////////////////////////////////////////////////////////////////////////
// Main loops and setup
///////////////////////////////////////////////////////////////////////////////

func main() {
	fmt.Println("start")

	// Set up logging to the file "debug.log"
	file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Screen initialization failed: %+v", err)
	}
	// Screen setup based on
	// https://github.com/gdamore/tcell/blob/master/_demos/boxes.go
	if err = s.Init(); err != nil {
		log.Fatalf("%v\n", err)
		os.Exit(1)
	}

	s.Clear()
	s.Show()

	for {
		PlayGame(s)
		wait := true
		for wait {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					s.Fini()
					return
				case tcell.KeyEnter:
					wait = false
				}
			}
		}
	}
}

// PlayGame has the main loop for the game of solitaire.
func PlayGame(s tcell.Screen) {
	var game Game = Deal()

	var gameWon bool = false

	// for loop based on https://github.com/gdamore/tcell/blob/master/_demos/boxes.go
	for {
		s.Clear()
		game.Render(s, 1, 1)
		s.Show()

		if gameWon {
			RenderGameWon(s, 1, 1)
			s.Show()
			return
		}

		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				s.Fini()
				return
			case tcell.KeyCtrlL:
				s.Sync()
			case tcell.KeyUp, tcell.KeyDown: // up and down do same thing
				game.Up()
			case tcell.KeyRight:
				game.Right()
			case tcell.KeyLeft:
				game.Left()
			case tcell.KeyEnter:
				gameWon = game.Select()
			case tcell.KeyRune:
				if ev.Rune() == ' ' {
					gameWon = game.Select()
				}
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// Game state modification functions
///////////////////////////////////////////////////////////////////////////////

// CreateDeck creates the deck with all cards. The Deck
// contains 2 full standard playing card decks (without
// jokers).
func CreateDeck() Deck {
	var deck Deck = NewDeck(52)
	// two full deck of cards
	for i := 0; i < 2; i++ {
		for s := Spades; s <= Diamonds; s++ {
			for v := Ace; v <= King; v++ {
				deck.Add(Card{s, v})
			}
		}
	}
	deck.Shuffle()
	return deck
}

// Deal creates all status needed to start the game,
// and returns it in the Game struct.
func Deal() Game {
	var game Game
	var deck Deck = CreateDeck()
	// first four piles get 6 cards
	for p := 0; p < 4; p++ {
		for c := 0; c < 5; c++ {
			var card Card = deck.Draw()
			game.piles[p].invisible.Add(card)
		}
	}
	// last for get 5 cards
	for p := 4; p < NUM_PILES; p++ {
		for c := 0; c < 4; c++ {
			var card Card = deck.Draw()
			game.piles[p].invisible.Add(card)
		}
	}
	// the top card is visible
	for p := 0; p < NUM_PILES; p++ {
		var card Card = deck.Draw()
		game.piles[p].visible.Add(card)
	}

	game.deck = deck
	game.highlighted.numCards = 1

	return game
}

// MoreCards deals another layer of cards onto the piles from the deck
func (game *Game) MoreCards() {
	for i := 0; i < NUM_PILES; i++ {
		game.piles[i].visible.Add(game.deck.Draw())
	}
}

// MoveCards attempts to move the selected cards to the highlighted pile.
func (game *Game) MoveCards() {
	Assert(game.highlighted.y == 1, "game.highlighted.y == 1")
	var topMovedCard Card = game.piles[game.selected.x].PeekNthCard(game.selected.numCards - 1)
	var cardMovedOnto Card = game.piles[game.highlighted.x].PeekNthCard(0)
	var canMove bool = (topMovedCard.value == cardMovedOnto.value-1) ||
		game.piles[game.highlighted.x].IsEmpty()
	if canMove {
		topNCards := game.piles[game.selected.x].GetTopNCards(game.selected.numCards)
		for _, v := range topNCards {
			game.piles[game.highlighted.x].visible.Add(v)
		}
	}
}

// IsFullStack returns true if the first 13 cards are a full stack
func IsFullStack(cards []Card) bool {
	log.Print("IsFullStack!", len(cards), NUM_VALUES)
	if len(cards) < NUM_VALUES {
		//fmt.Println("length: ", len(cards), NUM_VALUES)
		return false
	}

	stopVal := 0
	if len(cards)-NUM_VALUES > 0 {
		stopVal = len(cards) - NUM_VALUES
	}

	currValue := Ace
	for i := len(cards) - 1; i >= stopVal; i-- {

		if cards[i].value != currValue {
			return false
		}
		currValue++
	}
	return true
}

// CheckStacks looks for piles that contain a full stack of
// cards, and deletes off the full stacks.
func (game *Game) CheckStacks() {
	log.Print("Checking stacks!")
	for i := 0; i < NUM_PILES; i++ {
		if IsFullStack(game.piles[i].visible.cards) {
			game.piles[i].GetTopNCards(NUM_VALUES)
		}
	}
}

// CheckWon checks if there are no more cards and so the user
// has won. Returns true is the user has won, and false otherwise.
func (game *Game) CheckWon() bool {
	for i := 0; i < NUM_PILES; i++ {
		if !game.piles[i].IsEmpty() {
			return false
		}
	}
	// if theres nothing more in the deck or the piles, then we won
	return game.deck.IsEmpty()
}

///////////////////////////////////////////////////////////////////////////////
// Player move functions
///////////////////////////////////////////////////////////////////////////////

// Up makes changes for the user pressing the up or down arrows.
// Moves between deck and piles.
func (game *Game) Up() {
	game.highlighted.y = (game.highlighted.y + 1) % 2
	game.highlighted.numCards = 1
}

// Left makes changes for the user pressing the left arrow.
// Moves highlighted cursor one to the left.
func (game *Game) Left() {
	if game.highlighted.x == 0 {
		game.highlighted.x = NUM_PILES - 1
	} else {
		game.highlighted.x = (game.highlighted.x - 1) % NUM_PILES
	}
	game.highlighted.numCards = 1
}

// Right makes the changes for the user pressing the right arrow.
// Moves highlighted cursor one to the right.
func (game *Game) Right() {
	game.highlighted.x = (game.highlighted.x + 1) % NUM_PILES
	game.highlighted.numCards = 1
}

// Select makes the changes for the user pressing Select.
// When the deck is highlighted, Select deals more cards.
// When a pile is highlighted but not selected, Select
// selects the top card.
// When a pile is highlighted and cards from it are
// selected, Select selects one more card from that pile,
// if allowed.
func (game *Game) Select() bool {
	if game.highlighted.y == 1 {
		// The user has a pile highlighted
		if !game.toMove {
			// if nothing is selected, select the first card of whatever
			// pile is highlighted
			game.toMove = true
			game.selected.x = game.highlighted.x
			game.selected.y = game.highlighted.y
			game.selected.numCards = game.highlighted.numCards
		} else {
			// if we already have something selected
			selectedIsHighlighted := (game.highlighted.y == 1 &&
				game.highlighted.x == game.selected.x)
			if selectedIsHighlighted {
				// if the pile that is currently selected is highlighted,
				// try to select one more card from that pile. You can only
				// select multiple cards together if they are all moveable
				// together.
				if game.piles[game.selected.x].TopNMovable(game.selected.numCards + 1) {
					game.selected.numCards++
				}
			} else {
				// Try to move the selected cards to the new pile
				game.MoveCards()
				game.toMove = false
				// // The user is trying to select a pile other than what
				// // has been selected, so we change the selection to be
				// // whatever the user has highlighted.
				// game.selected.x = game.highlighted.x
				// game.selected.y = game.highlighted.y
				// game.selected.numCards = game.highlighted.numCards
			}
		}

	} else {
		// the user has pressed enter while the deck is highlighted.
		// Get more cards from the deck.
		game.MoreCards()
	}

	// The player pressing enter can trigger any given pile to
	// now have a full stack.
	game.CheckStacks()
	return game.CheckWon()
}

///////////////////////////////////////////////////////////////////////////////
// Graphics
///////////////////////////////////////////////////////////////////////////////

// Render renders the full current game
func (game Game) Render(s tcell.Screen, x int, y int) {
	if !game.deck.IsEmpty() {
		// game.deck.cards[0].RenderFlipped(s, x, y, (game.highlighted.y == 0))
		game.deck.cards[0].RenderFlipped(s, x, y)
	}
	for i := 0; i < NUM_PILES; i++ {
		game.piles[i].Render(s, x+(CARD_WIDTH+2)*i, y+CARD_HEIGHT+2)
	}
	style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

	hglt := game.highlighted
	var higBoxX int = 1
	var higBoxY int = 1
	if hglt.y == 1 {
		higBoxX = x + hglt.x*(CARD_WIDTH+2)
		distFromPileTop := game.piles[hglt.x].Height() - hglt.numCards*2
		higBoxY = y + CARD_HEIGHT + 2 + distFromPileTop
	}
	var higBox Box = Box{s, higBoxX, higBoxY,
		higBoxX + CARD_WIDTH, higBoxY + CARD_HEIGHT + ((hglt.numCards - 1) * 2),
		style, "", true}
	higBox.Draw()

	if game.toMove {
		sel := game.selected
		style = tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorYellow)
		var selBoxX int = 1
		var selBoxY int = 1
		if sel.y == 1 {
			selBoxX = x + sel.x*(CARD_WIDTH+2)
			distFromPileTop := game.piles[sel.x].Height() - (sel.numCards * 2)
			selBoxY = y + CARD_HEIGHT + 2 + distFromPileTop
		}
		var selBox Box = Box{s, selBoxX, selBoxY,
			selBoxX + CARD_WIDTH, selBoxY + CARD_HEIGHT + ((sel.numCards - 1) * 2),
			style, "", true}
		selBox.Draw()
	}
}

// RenderGameWon renders the message that the game was won.
func RenderGameWon(s tcell.Screen, x int, y int) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen)
	var box Box = Box{s, x, y,
		x + 10*(CARD_WIDTH+1), y + 2*CARD_HEIGHT + 2,
		style, "You won! Press ESC to leave, and enter to restart", false}
	box.Draw()
}

///////////////////////////////////////////////////////////////////////////////
// Utilities
///////////////////////////////////////////////////////////////////////////////

func Assert(b bool, s string) {
	if !b {
		log.Fatalf(s)
	}
}
