// This file is part of Ringot.
/*
Copyright 2016 tSU-RooT <tsu.root@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"unicode/utf8"
)

const (
	inputMode   = "*Input Mode*"
	commandMode = "*Command Mode*"
	confirmText = "ok?[Enter/C-g]"
	Margin      = 18
)

type buffer struct {
	content      []byte
	cursorX      int
	cursorOffset int
	mode         string
	process      func(string)
	inputing     bool
	confirm      bool
	confirmLock  lock
	commanding   bool

	linePosInfo int
}

func newBuffer() *buffer {
	return &buffer{}
}

func (bf *buffer) draw() {
	t := bf.mode
	if bf.inputing {
		if bf.commanding {
			t = commandMode
		} else {
			t = inputMode
		}

	}
	// Draw upper line
	width, height := getTermSize()
	fillLine(0, height-2, ColorGray2)

	info := fmt.Sprintf("User:@%s [%s]", user.ScreenName, user.UserName)
	x := width - runewidth.StringWidth(info) - 1
	drawText(info, x, height-2, ColorGreen, ColorGray2)

	info = fmt.Sprintf("L%d", bf.linePosInfo)
	x -= runewidth.StringWidth(info) + 1
	drawText(info, x, height-2, ColorWhite, ColorGray2)

	clength := utf8.RuneCountInString(string(bf.content))
	if bf.inputing && !bf.commanding && clength >= 20 {
		info = fmt.Sprintf("length:(%d)", clength)
		x -= runewidth.StringWidth(info) + 1
		drawText(info, x, height-2, ColorWhite, ColorGray2)
	}

	x = 2
	drawText(t, x, height-2, ColorYellow, ColorGray2)
	x += runewidth.StringWidth(t)
	termbox.SetCell(x, height-2, ' ', ColorBackground, ColorGray2)
	x++

	// Draw lower line
	con := string(bf.content[bf.cursorOffset:])
	if bf.commanding {
		con = ":" + con
	}
	drawText(con, 0, height-1, ColorWhite, ColorBackground)
	x = runewidth.StringWidth(con)
	if bf.confirm {
		x++
		t := confirmText
		drawText(t, x, height-1, ColorRed, ColorBackground)
		x += runewidth.StringWidth(t)
	}
	fillLine(x, height-1, ColorBackground)
}

func (bf *buffer) runeUnderCursor() (rune, int) {
	return utf8.DecodeRune(bf.content[bf.cursorX:])
}

func (bf *buffer) insertRune(r rune) {
	var u [utf8.UTFMax]byte
	s := utf8.EncodeRune(u[:], r)
	bf.content = byteSliceInsert(bf.content, u[:s], bf.cursorX)
	bf.cursorMoveForward()
}

func (bf *buffer) deleteRuneBackward() {
	if bf.cursorX == 0 {
		return
	}
	bf.cursorMoveBackward()
	_, s := bf.runeUnderCursor()
	bf.content = byteSliceRemove(bf.content, bf.cursorX, bf.cursorX+s)
}

func (bf *buffer) cursorMoveBackward() {
	if bf.cursorX <= 0 {
		return
	}
	_, s := utf8.DecodeLastRune(bf.content[:bf.cursorX])
	bf.cursorX -= s
	if bf.cursorOffset > 0 && bf.cursorX <= bf.cursorOffset+Margin {
		_, s := utf8.DecodeLastRune(bf.content[:bf.cursorOffset])
		bf.cursorOffset -= s
	}
	return
}

func (bf *buffer) cursorMoveForward() {
	if bf.cursorX >= len(bf.content) {
		return
	}
	_, size := utf8.DecodeRune(bf.content[bf.cursorX:])
	bf.cursorX += size
	width, _ := getTermSize()
	cw := runewidth.StringWidth(string(bf.content[bf.cursorOffset:bf.cursorX]))
	if cw >= width-Margin {
		_, size = utf8.DecodeRune(bf.content[bf.cursorOffset:])
		bf.cursorOffset += size
	}
}

func (bf *buffer) cursorMoveToTop() {
	bf.cursorX = 0
	bf.cursorOffset = 0
}

func (bf *buffer) cursorMoveToBottom() {
	width, _ := getTermSize()
	bf.cursorX = len(bf.content)
	cw := runewidth.StringWidth(string(bf.content))
	if cw >= width-Margin {
		bf.cursorOffset = cw - (width - Margin)
	}

}

func (bf *buffer) updateCursorPosition() {
	if (!bf.inputing) || bf.confirm || bf.cursorX < 0 || bf.cursorX > len(bf.content) {
		termbox.HideCursor()
		return
	}
	_, h := getTermSize()
	x := runewidth.StringWidth(string(bf.content[bf.cursorOffset:bf.cursorX]))
	if bf.commanding {
		x++
	}
	termbox.SetCursor(x, h-1)
}

func (bf *buffer) setContent(s string) {
	b := []byte(s)
	bf.content = make([]byte, len(b), 180)
	copy(bf.content, b)
	bf.cursorX = len(b)
	bf.cursorOffset = 0
}

func (bf *buffer) setModeStr(m viewmode) {
	s := ""
	switch m {
	case home:
		s = "*Timeline View*"
	case mention:
		s = "*Mention View*"
	case conversation:
		s = "*Conversation View*"
	case usertimeline:
		s = "*UserTimeline View*"
	case list:
		s = "*List View*"
	}
	bf.mode = s
}
