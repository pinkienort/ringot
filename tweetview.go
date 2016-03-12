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
	"github.com/ChimeraCoder/anaconda"
	"github.com/mattn/go-runewidth"
	"strings"
	"time"
)

type tweetview struct {
	tweets         []tweetstatus
	cursorPosition int
	scroll         int
	scrollOffset   int
}

func newTweetview() *tweetview {
	return &tweetview{
		tweets: []tweetstatus{tweetstatus{ReloadMark: true}},
	}
}

func (tv *tweetview) isEmpty() bool {
	if len(tv.tweets) == 0 {
		return true
	} else if (tv.tweets[0].ReloadMark || tv.tweets[0].Empty) && len(tv.tweets) == 1 {
		return true
	}
	return false
}

func (tv *tweetview) cursorDown() {
	_, h := getTermSize()
	if tv.cursorPosition+1 < len(tv.tweets) {
		tv.cursorPosition++
		sum := 0
		tweets := tv.tweets[:tv.cursorPosition+1]
		for _, t := range tweets {
			sum += t.countLines()
		}
		sub := (sum - (tv.scroll - tv.scrollOffset)) - (h - 2)
		if sub >= 0 {
			tv.scroll += sub
		}
	}
}

func (tv *tweetview) cursorUp() {
	if tv.cursorPosition > 0 {
		tv.cursorPosition--
		sum := 0
		tweets := tv.tweets[:tv.cursorPosition]
		for _, t := range tweets {
			sum += t.countLines()
		}
		sub := sum - tv.scroll
		if sub <= 0 {
			tv.scroll += sub
			if tv.scroll < 0 {
				tv.scroll = 0
			}
		}
	}
}

func (tv *tweetview) addNewTweet(tss []tweetstatus) {
	if len(tv.tweets) > 1 {
		tv.scroll += sumTweetLines(tss)
		tv.cursorPosition += len(tss)
	}
	tv.tweets = append(tss, tv.tweets...)
}

func (tv *tweetview) addIntervalTweet(tss []tweetstatus) {
	adjustScroll := false
	if tv.tweets[tv.cursorPosition].ReloadMark {
		adjustScroll = true
	}

	i := 0
	for i = range tv.tweets {
		if tv.tweets[i].ReloadMark {
			break
		}
	}
	tweets := tv.tweets[:i]
	if tweets[len(tweets)-1].Content.Id == tss[0].Content.Id {
		tss = tss[1:]
	}
	t := append(tweets, tss...)
	t = append(t, tweetstatus{ReloadMark: true})
	tv.tweets = t
	if adjustScroll {
		tv.cursorPosition--
		tv.cursorDown()
	}
}

func (tv *tweetview) draw() {
	width, _ := getTermSize()
	y := -(tv.scroll - tv.scrollOffset)

	index := 0
	now := time.Now()

	tweets := tv.tweets

	for ; index < len(tweets); index++ {
		tweetstatus := tweets[index]
		countLine := tweetstatus.countLines()
		if y+countLine < 0 {
			y += countLine
			continue
		}
		bgColor := ColorBackground
		cursorColor := ColorBackground
		selected := index == tv.cursorPosition
		if selected {
			cursorColor = ColorGray3
		}

		if tweetstatus.ReloadMark {
			if selected {
				drawText(" ⟳ Reload", 0, y, ColorWhite, ColorGray1)
			} else {
				drawText(" ⟳ Reload", 0, y, ColorWhite, bgColor)
			}
			y++
			continue
		}
		tweet := tweetstatus.Content
		favorited := tweet.Favorited
		retweeted := tweet.Retweeted
		retweetedBy := ""
		labelColor := generateLabelColorByUserID(tweet.User.Id)
		if tweet.RetweetedStatus != nil {
			retweetedBy = tweet.User.ScreenName
			tweet = tweet.RetweetedStatus
			labelColor = ColorPink
		}
		text := tweet.Text

		x := 0
		drawText(" ", x, y, ColorBackground, labelColor)
		// Draw Name
		x = 1
		drawText(" ", x, y, ColorBackground, cursorColor)
		x = 2
		drawText("@"+tweet.User.ScreenName, x, y, labelColor, bgColor)
		x += runewidth.StringWidth("@"+tweet.User.ScreenName) + 1
		drawText(tweet.User.Name, x, y, ColorWhite, bgColor)
		x += runewidth.StringWidth(tweet.User.Name) + 1
		if favorited {
			drawText("★", x, y, ColorYellow, bgColor)
			x += runewidth.StringWidth("★") + 1
		}
		if retweeted {
			drawText("RT", x, y, ColorGreen, ColorLowlight)
			x += runewidth.StringWidth("RT") + 1
		}
		if retweetedBy != "" {
			drawText("ReTweeted By @"+retweetedBy, x, y, ColorRed, bgColor)
		}
		y++

		lines := strings.Split(runewidth.Wrap(text, width-2), "\n")
		for i, t := range lines {
			drawText(" ", 0, y+i, ColorBackground, labelColor)
			drawText(" ", 1, y+i, ColorWhite, cursorColor)
			drawTextWithAutoNotice(t, 2, y+i, ColorWhite, bgColor)
		}
		y += len(lines)

		// Draw Tweet Detail
		createdAtTime, err := tweet.CreatedAtTime()
		if err != nil {
			continue
		}
		sub := now.Sub(createdAtTime)
		var strTime string
		if sub <= time.Second*30 {
			strTime = "now"
		} else if sub <= time.Minute*5 {
			strTime = "A few minutes ago"
		} else if sub <= time.Hour*2 {
			m := sub / time.Minute
			strTime = fmt.Sprintf("%d minutes ago", m)
		} else if sub <= time.Hour*36 {
			h := sub / time.Hour
			strTime = fmt.Sprintf("%d hours ago", h)
		} else {
			strTime = fmt.Sprintf("%d/%d/%d", createdAtTime.Year(), createdAtTime.Month(), createdAtTime.Day())
		}

		drawText(" ", 0, y, ColorBackground, labelColor)
		x = 1
		drawText(" "+strTime, x, y, ColorGray3, cursorColor)
		x = 2
		drawText(strTime, x, y, ColorGray1, bgColor)
		x += 1 + runewidth.StringWidth(strTime)
		if tweet.RetweetCount > 0 {
			strRT := fmt.Sprintf("RT %d", tweet.RetweetCount)
			drawText(" "+strRT, x, y, ColorGreen, bgColor)
			x += 1 + runewidth.StringWidth(strRT)
		}
		if tweet.FavoriteCount > 0 {
			strFav := fmt.Sprintf("Fav %d", tweet.FavoriteCount)
			drawText(" "+strFav, x, y, ColorYellow, bgColor)
			x += 1 + runewidth.StringWidth(strFav)
		}

		y++

	}
}

func (tv *tweetview) resetScroll() {
	sum := 0
	_, h := getTermSize()
	var tweets []tweetstatus
	tweets = tv.tweets[:tv.cursorPosition]

	for _, t := range tweets {
		sum += t.countLines()
	}
	sub := sum - tv.scroll
	if sub <= 0 {
		tv.scroll = sum
	} else {
		cl := tv.tweets[tv.cursorPosition].countLines()
		sub = (sum + cl - (tv.scroll - tv.scrollOffset)) - (h - 2)
		if sub >= 0 {
			tv.scroll = sum + -(h - 2 - cl)
		}
	}

}

type tweetstatus struct {
	Content    *anaconda.Tweet
	ReloadMark bool
	Empty      bool
}

func (status *tweetstatus) countLines() int {
	if status.Empty {
		return 0
	} else if status.ReloadMark {
		return 1
	}
	w, _ := getTermSize()
	tweet := status.Content
	if tweet.RetweetedStatus != nil {
		tweet = tweet.RetweetedStatus
	}
	text := tweet.Text
	lines := strings.Split(runewidth.Wrap(text, w-2), "\n")
	lineCount := 1 + len(lines) + 1
	return lineCount
}

func (status *tweetstatus) isFavorited() bool {
	if status.Content != nil {
		return status.Content.Favorited
	}
	return false
}
func (status *tweetstatus) setFavorited(b bool) {
	if status.Content != nil {
		status.Content.Favorited = b
	}
}

func (status *tweetstatus) isRetweeted() bool {
	if status.Content != nil {
		return status.Content.Retweeted
	}
	return false
}

func (status *tweetstatus) setRetweeted(b bool) {
	if status.Content != nil {
		status.Content.Retweeted = b
	}
}
