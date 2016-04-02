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
	"strings"
	"unicode/utf8"
)

const (
	inputMode   = "*Input Mode*"
	commandMode = "*Command Mode*"
	confirmText = "ok?[Enter/C-g]"
)

type buffer struct {
	content     []byte
	state       string
	cursorX     int
	mode        string
	process     func(string)
	inputing    bool
	confirm     bool
	confirmLock lock
	commanding  bool

	linePosInfo int
}

func newBuffer() *buffer {
	return &buffer{}
}

func (bf *buffer) draw() {
	if bf.inputing && !bf.commanding {
		bf.drawTweetInputArea()
	} else {
		bf.drawCommandInputField()
	}

}

func (bf *buffer) drawTweetInputArea() {
	width, height := getTermSize()
	// Draw upper line
	fillLine(0, height-5, ColorGray2)

	info := fmt.Sprintf("User:@%s [%s]", user.ScreenName, user.UserName)
	x := width - runewidth.StringWidth(info) - 1
	drawText(info, x, height-5, ColorGreen, ColorGray2)

	info = fmt.Sprintf("L%d", bf.linePosInfo)
	x -= runewidth.StringWidth(info) + 1
	drawText(info, x, height-5, ColorWhite, ColorGray2)

	clength := utf8.RuneCountInString(string(bf.content))
	if bf.inputing && !bf.commanding && clength >= 20 {
		info = fmt.Sprintf("length:(%d)", clength)
		x -= runewidth.StringWidth(info) + 1
		drawText(info, x, height-5, ColorWhite, ColorGray2)
	}

	x = 2
	drawText(inputMode, x, height-5, ColorYellow, ColorGray2)
	x += runewidth.StringWidth(inputMode) + 1
	if bf.confirm {
		drawText(confirmText, x, height-5, ColorRed, ColorGray2)
	}

	// Draw Input Area
	x = 0
	text := string(bf.content)
	lines := strings.Split(runewidth.Wrap(text, width), "\n")

	for i := 0; i < 4; i++ {
		if i >= len(lines) {
			fillLine(0, height-4+i, ColorBackground)
			continue
		}
		drawText(lines[i], 0, height-4+i, ColorWhite, ColorBackground)
		x = runewidth.StringWidth(lines[i])
		fillLine(x, height-4+i, ColorBackground)
	}

}

func (bf *buffer) drawCommandInputField() {
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
	x = 0
	if bf.inputing {
		con := string(bf.content)
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
	} else {
		drawText(bf.state, x, height-1, ColorWhite, ColorBackground)
		x += runewidth.StringWidth(bf.state)
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

func (bf *buffer) insertLF() {
	w, _ := getTermSize()
	text := string(bf.content)
	lines := strings.Split(runewidth.Wrap(text, w), "\n")
	if len(lines) < 4 {
		bf.insertRune('\n')
	}
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
}

func (bf *buffer) cursorMoveForward() {
	if bf.cursorX >= len(bf.content) {
		return
	}
	_, size := utf8.DecodeRune(bf.content[bf.cursorX:])
	bf.cursorX += size
}

var (
	_, LFByteSize = utf8.DecodeRune([]byte("\n"))
)

func (bf *buffer) cursorMoveUp() {
	if bf.cursorX > len(bf.content) || bf.commanding {
		return
	}
	lines := strings.Split(string(bf.content[:bf.cursorX]), "\n")
	if len(lines) <= 1 {
		return
	}
	w := runewidth.StringWidth(lines[len(lines)-1])
	if runewidth.StringWidth(lines[len(lines)-2]) <= w {
		x := 0
		for i := 0; i < len(lines)-1; i++ {
			t := []byte(lines[i])
			for len(t) > 0 {
				_, s := utf8.DecodeRune(t)
				x += s
				t = t[s:]
			}
			x += LFByteSize
		}
		// Don't count the last LF
		x -= LFByteSize
		bf.cursorX = x
		return
	}

	x := 0
	for i := 0; i < len(lines)-2; i++ {
		t := []byte(lines[i])
		for len(t) > 0 {
			_, s := utf8.DecodeRune(t)
			x += s
			t = t[s:]
		}
		x += LFByteSize
	}
	t := []byte(lines[len(lines)-2])
	wc := 0
	for len(t) > 0 {
		r, s := utf8.DecodeRune(t)
		wc += runewidth.RuneWidth(r)
		if wc > w {
			break
		}
		x += s
		t = t[s:]
	}
	bf.cursorX = x
}

func (bf *buffer) cursorMoveDown() {
	if bf.cursorX > len(bf.content) || bf.commanding {
		return
	}
	lines1 := strings.Split(string(bf.content[:bf.cursorX]), "\n")
	lines2 := strings.Split(string(bf.content), "\n")
	if len(lines1) < 1 || len(lines2) <= len(lines1) {
		return
	}
	width := runewidth.StringWidth(lines1[len(lines1)-1])
	l := len(lines1)
	if runewidth.StringWidth(lines2[l]) <= width {
		x := 0
		for i := 0; i < len(lines2); i++ {
			t := []byte(lines2[i])
			for len(t) > 0 {
				_, s := utf8.DecodeRune(t)
				x += s
				t = t[s:]
			}
			x += LFByteSize
		}
		// Don't count the last LF
		x -= LFByteSize
		bf.cursorX = x
		return
	}

	x := 0
	for i := 0; i < l; i++ {
		t := []byte(lines2[i])
		for len(t) > 0 {
			_, s := utf8.DecodeRune(t)
			x += s
			t = t[s:]
		}
		x += LFByteSize
	}
	t := []byte(lines2[l])
	wc := 0
	for len(t) > 0 {
		r, s := utf8.DecodeRune(t)
		wc += runewidth.RuneWidth(r)
		if wc > width {
			break
		}
		x += s
		t = t[s:]
	}
	bf.cursorX = x

}

func (bf *buffer) cursorMoveToLineTop() {
	if bf.commanding || len(bf.content) == 0 {
		bf.cursorX = 0
	} else {
		lines := strings.Split(string(bf.content[:bf.cursorX]), "\n")
		if len(lines) <= 1 {
			bf.cursorX = 0
			return
		}
		x := 0
		for i := 0; i < len(lines)-1; i++ {
			t := []byte(lines[i])
			for len(t) > 0 {
				_, s := utf8.DecodeRune(t)
				x += s
				t = t[s:]
			}
			x += LFByteSize
		}
		bf.cursorX = x
	}
}

func (bf *buffer) cursorMoveToLineBottom() {
	if bf.commanding {
		bf.cursorX = len(bf.content)
	} else {
		lines := strings.Split(string(bf.content), "\n")
		if len(lines) <= 1 {
			bf.cursorX = len(bf.content)
			return
		}
		x := 0
		for i := 0; i < len(lines); i++ {
			b := false
			t := []byte(lines[i])
			for len(t) > 0 {
				_, s := utf8.DecodeRune(t)
				x += s
				if x >= bf.cursorX {
					b = true
				}
				t = t[s:]
			}
			if b {
				break
			}
			x += LFByteSize
		}
		bf.cursorX = x
	}
}

func (bf *buffer) updateCursorPosition() {
	if (!bf.inputing) || bf.confirm || bf.cursorX < 0 || bf.cursorX > len(bf.content) {
		termbox.HideCursor()
		return
	}
	w, h := getTermSize()
	if bf.commanding {
		x := runewidth.StringWidth(string(bf.content[:bf.cursorX]))
		if bf.commanding {
			x++
		}
		termbox.SetCursor(x, h-1)
	} else {
		text := string(bf.content[:bf.cursorX])
		lines := strings.Split(runewidth.Wrap(text, w), "\n")
		x := runewidth.StringWidth(lines[len(lines)-1])
		if x == w {
			termbox.SetCursor(0, h-4+len(lines))
		} else {
			termbox.SetCursor(x, h-5+len(lines))
		}
	}
}

func (bf *buffer) setContent(s string) {
	b := []byte(s)
	bf.content = make([]byte, len(b), 180)
	copy(bf.content, b)
	bf.cursorX = len(b)
}

func (bf *buffer) setState(s string) {
	bf.state = s
}

func (bf *buffer) clear() {
	bf.setContent("")
	bf.setState("")
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
