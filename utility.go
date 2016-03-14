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
	"errors"
	"github.com/ChimeraCoder/anaconda"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"
)

func byteSliceRemove(bytes []byte, from int, to int) []byte {
	copy(bytes[from:], bytes[to:])
	return bytes[:len(bytes)+from-to]
}

func byteSliceInsert(dst []byte, src []byte, pos int) []byte {
	length := len(dst) + len(src)
	if cap(dst) < length {
		s := make([]byte, len(dst), length)
		copy(s, dst)
		dst = s
	}
	dst = dst[:length]
	copy(dst[pos+len(src):], dst[pos:])
	copy(dst[pos:], src)
	return dst

}

func centeringStr(str string, width int) string {
	sub := width - len(str)
	if sub <= 0 {
		return str
	}
	val := ""
	if sub%2 == 0 {
		for i := 0; i < (sub / 2); i++ {
			val += " "
		}
	} else {
		for i := 0; i < (sub/2)+1; i++ {
			val += " "
		}
	}
	val += str

	for i := 0; i < (sub / 2); i++ {
		val += " "
	}
	return val
}

func drawText(str string, x int, y int, fg termbox.Attribute, bg termbox.Attribute) {
	i := 0
	for _, c := range str {
		termbox.SetCell(x+i, y, c, fg, bg)
		i += runewidth.RuneWidth(c)
	}
}

func drawTextWithAutoNotice(str string, x int, y int, fg termbox.Attribute, bg termbox.Attribute) {
	pos := 0
	foreColor := fg
	backColor := bg
	fgChanging := false
	bgChanging := false
	t := []byte(str)
	for {
		if len(t) == 0 {
			break
		}
		c, s := utf8.DecodeRune(t)
		if !(bgChanging || fgChanging) && len(t) > s {
			if c == '@' {
				tc, _ := utf8.DecodeRune(t[s:])
				if isScreenNameUsable(tc) {
					backColor = ColorLowlight
					bgChanging = true
				}
			} else {
				found := false
				s2 := 0
				if c == ' ' {
					var tc rune
					tc, s2 = utf8.DecodeRune(t[s:])
					if tc == '#' {
						found = true
					}
				} else if c == '#' && pos == 0 {
					found = true
				}
				if found {
					tc, _ := utf8.DecodeRune(t[s+s2:])
					if tc != ' ' {
						foreColor = ColorBlue
						fgChanging = true
					}
				}
			}
		} else {
			if bgChanging && !isScreenNameUsable(c) {
				backColor = bg
				bgChanging = false
			} else if fgChanging && c == ' ' {
				tc, _ := utf8.DecodeRune(t[s:])
				if tc != '#' {
					foreColor = fg
					fgChanging = false
				}
			}
		}

		termbox.SetCell(x+pos, y, c, foreColor, backColor)
		pos += runewidth.RuneWidth(c)
		t = t[s:]
	}
}

func isScreenNameUsable(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	} else if r >= 'A' && r <= 'Z' {
		return true
	} else if r >= '0' && r <= '9' {
		return true
	} else if r == '_' {
		return true
	}
	return false
}

func fillLine(offset int, y int, bg termbox.Attribute) {
	width, _ := getTermSize()
	x := offset
	for {
		if x >= width {
			break
		}
		termbox.SetCell(x, y, ' ', ColorBackground, bg)
		x++
	}
}

func generateLabelColorByUserID(id int64) termbox.Attribute {
	if val, ok := LabelColorMap[id]; ok {
		return LabelColors[val]
	}

	rand.Seed(id)
	val := rand.Intn(len(LabelColors))
	LabelColorMap[id] = val
	return LabelColors[val]
}

var (
	replacer = strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">")
)

func wrapTweets(tweets []anaconda.Tweet) []tweetstatus {
	result := make([]tweetstatus, len(tweets))
	for i := 0; i < len(tweets); i++ {
		tweets[i].Text = replacer.Replace(tweets[i].Text)
		for _, url := range tweets[i].Entities.Urls {
			tweets[i].Text = strings.Replace(tweets[i].Text, url.Url, url.Display_url, -1)
		}
		for _, media := range tweets[i].ExtendedEntities.Media {
			tweets[i].Text = strings.Replace(tweets[i].Text, media.Url, media.Display_url, -1)
		}
		result[i] = tweetstatus{Content: &tweets[i]}
	}
	return result
}

func sumTweetLines(tweetsStatusSlice []tweetstatus) int {
	sum := 0
	tweets := tweetsStatusSlice
	for _, t := range tweets {
		sum += t.countLines()
	}
	return sum
}

func openCommand(path string) {
	var commandName string
	switch runtime.GOOS {
	case "linux":
		commandName = "xdg-open"
	case "darwin":
		commandName = "open"
	default:
		return

	}
	exec.Command(commandName, path).Run()
}

const (
	tempDir = "ringot"
)

func downloadMedia(url string) (fullpath string, err error) {
	_, filename := path.Split(url)
	fullpath = filepath.Join(os.TempDir(), tempDir, filename)
	if _, err := os.Stat(fullpath); err == nil {
		return "", os.ErrExist
	}

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		errors.New(res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	tempdir := filepath.Join(os.TempDir(), tempDir)
	if _, err := os.Stat(tempdir); err != nil {
		err := os.Mkdir(tempdir, 0775)
		if err != nil {
			return "", err
		}
	}
	file, err := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0664)
	if err != nil {
		return "", err
	}
	defer file.Close()
	file.Write(body)
	return fullpath, nil
}

func openMedia(url string) {
	fullpath, err := downloadMedia(url)
	if err != nil && err != os.ErrExist {
		panic(err)
		return
	}
	openCommand(fullpath)
}

func favoriteTweet(id int64) {
	_, err := api.Favorite(id)
	if err != nil {
		changeBufferState("Err:Favorite")
		return
	}
}

func unfavoriteTweet(id int64) {
	_, err := api.Unfavorite(id)
	if err != nil {
		changeBufferState("Err:Unfavorite")
		return
	}
}

func retweet(id int64) {
	_, err := api.Retweet(id, false)
	if err != nil {
		changeBufferState("Err:Retweet")
		return
	}
}

func changeBufferState(state string) {
	go func() { stateCh <- state }()
}

func getTermSize() (int, int) {
	return termWidth, termHeight
}

func setTermSize(w, h int) {
	termWidth, termHeight = w, h
}

type lock struct {
	mutex   sync.Mutex
	locking uint32
}

// Errors
var (
	ErrAlreayLocking = errors.New("already locking")
)

func (l *lock) lock() error {
	if atomic.LoadUint32(&l.locking) == 1 {
		return ErrAlreayLocking
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.locking == 0 {
		atomic.StoreUint32(&l.locking, 1)
	}
	return nil
}

func (l *lock) unlock() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.locking == 1 {
		atomic.StoreUint32(&l.locking, 0)
	}
}

func (l *lock) isLocking() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return atomic.LoadUint32(&l.locking) == 1
}
