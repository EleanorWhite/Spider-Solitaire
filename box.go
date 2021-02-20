package main

import "github.com/gdamore/tcell"

// Box is a box
type Box struct {
	s     tcell.Screen
	x1    int
	y1    int
	x2    int
	y2    int
	style tcell.Style
	text  string
	isEmpty bool
}

// Draw draws a box
func (b Box) Draw() {
	if b.isEmpty {
		drawEmptyBox(b.s, b.x1, b.y1, b.x2, b.y2, b.style, b.text)
	} else {
		drawBox(b.s, b.x1, b.y1, b.x2, b.y2, b.style, b.text)
	}
}


///////////////////////////////////////////////////////////////////////////////
// Graphics Functions
///////////////////////////////////////////////////////////////////////////////
// Graphics functions taken and modified from 
// https://github.com/gdamore/tcell/blob/master/_demos/mouse.go

// drawEmptyBox draws the outline of a box.
func drawEmptyBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
}

// emitStr puts string str in between the coordinates x1, y1 and x2, y2. 
// Will truncate the str if it does not fit. 
func emitStr(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, str string) {
	x := x1
	y := y1
	for _, c := range str {
		s.SetContent(x, y, c, nil, style)
		newX := (((x-x1)+1) % (y2-y1)) + x1
		if newX <= x {
			y++
			if(y > y2) {
				return
			}
		}
		x = newX
	}
}

// drawBox draws a filled in box containing the string str.
func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, str string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}
	emitStr(s, x1 + 1, y1 + 1, x2 - 1, y2 - 1, style, str)
}

