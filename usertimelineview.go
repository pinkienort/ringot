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
	"net/url"
	"strconv"
	"strings"
)

type usertimelineview struct {
	*tweetview
	screenName string

	loading             lock
	loadNewTweetCh      chan []anaconda.Tweet
	loadIntervalTweetCh chan []anaconda.Tweet
}

func newUsertimelineview() *usertimelineview {
	return &usertimelineview{
		tweetview:           newTweetview(),
		loadNewTweetCh:      make(chan []anaconda.Tweet),
		loadIntervalTweetCh: make(chan []anaconda.Tweet),
	}
}

func (uv *usertimelineview) setUserScreenName(name string) {
	if uv.screenName != name {
		uv.screenName = name
		uv.tweets = []tweetstatus{tweetstatus{ReloadMark: true}}
	}
}

func (uv *usertimelineview) loadTweet(sinceID int64) {
	if uv.loading.isLocking() {
		return
	}
	uv.loading.lock()
	defer uv.loading.unlock()
	changeBufferState("Loading...")
	val := url.Values{}
	val.Add("count", strconv.Itoa(20))
	val.Add("screen_name", uv.screenName)
	if sinceID > 0 {
		val.Add("since_id", strconv.FormatInt(sinceID, 10))
	}
	timeline, err := api.GetUserTimeline(val)
	if err != nil {
		changeBufferState("Err:Loading")
		return
	}
	changeBufferState(fmt.Sprintf("Load!(%d tweets)", len(timeline)))
	uv.loadNewTweetCh <- timeline
}

func (uv *usertimelineview) loadIntervalTweet(maxID int64) {
	if uv.loading.isLocking() {
		return
	}
	uv.loading.lock()
	defer uv.loading.unlock()
	changeBufferState("Loading...")
	val := url.Values{}
	val.Add("count", strconv.Itoa(CountTweet))
	val.Add("screen_name", uv.screenName)
	if maxID > 0 {
		val.Add("max_id", strconv.FormatInt(maxID, 10))
	}
	timeline, err := api.GetUserTimeline(val)
	if err != nil {
		changeBufferState("Err:Loading")
		return
	}
	changeBufferState(fmt.Sprintf("Load!(%d tweets)", len(timeline)))
	uv.loadIntervalTweetCh <- timeline
}

func (uv *usertimelineview) draw() {
	tweets := uv.tweets
	if len(tweets) == 0 || tweets[0].ReloadMark {
		drawText("Now Loading...", 0, 0, ColorWhite, ColorBackground)
		return
	}

	width, _ := getTermSize()
	user := tweets[0].Content.User
	lines := strings.Split(runewidth.Wrap(user.Description, width), "\n")

	// slide for Users profile space
	uv.scrollOffset = (3 + len(lines))
	uv.tweetview.draw()

	x := 0
	y := 0
	text := fmt.Sprintf("@%s", user.ScreenName)
	labelColor := generateLabelColorByUserID(user.Id)
	drawText(text, 0, y, labelColor, ColorGray2)
	x += runewidth.StringWidth(text)
	drawText(" ", x, y, ColorWhite, ColorGray2)
	x++
	text = fmt.Sprintf("%s", user.Name)
	drawText(text, x, y, ColorWhite, ColorGray2)
	x += runewidth.StringWidth(text)
	if user.Protected {
		text = "[Protected]"
		drawText(text, x, y, ColorWhite, ColorGray2)
		x += runewidth.StringWidth(text)
	}
	fillLine(x, y, ColorGray2)
	y++

	for _, t := range lines {
		drawText(t, 0, y, ColorWhite, ColorGray2)
		x = runewidth.StringWidth(t)
		fillLine(x, y, ColorGray2)
		y++
	}
	x = 0
	if user.URL != "" {
		text = "URL:" + user.URL
	} else {
		text = "URL:None"
	}
	drawText(text, x, y, ColorWhite, ColorGray2)
	x += runewidth.StringWidth(text)
	fillLine(x, y, ColorGray2)
	y++

	var ws int
	if width%4 == 0 {
		ws = width / 4
	} else {
		ws = (width / 4) + 1
	}
	x = 0
	text = centeringStr(fmt.Sprintf("Tweets:%d", user.StatusesCount), ws)
	drawText(text, x, y, ColorWhite, ColorBlue)
	x += runewidth.StringWidth(text)
	text = centeringStr(fmt.Sprintf("Follwing:%d", user.FriendsCount), ws)
	drawText(text, x, y, ColorWhite, ColorRed)
	x += runewidth.StringWidth(text)
	text = centeringStr(fmt.Sprintf("Follower:%d", user.FollowersCount), ws)
	drawText(text, x, y, ColorWhite, ColorGreen)
	x += runewidth.StringWidth(text)
	text = centeringStr(fmt.Sprintf("Favorite:%d", user.FavouritesCount), ws)
	drawText(text, x, y, ColorWhite, ColorYellow)

}

func (uv *usertimelineview) resetScroll() {
	if len(uv.tweets) == 0 || uv.tweets[0].Content == nil {
		return
	}
	width, _ := getTermSize()
	user := uv.tweets[0].Content.User
	lines := strings.Split(runewidth.Wrap(user.Description, width), "\n")
	uv.scrollOffset = (3 + len(lines))

	uv.tweetview.resetScroll()
}
