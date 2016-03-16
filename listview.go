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
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/mattn/go-runewidth"
	"net/url"
	"strconv"
	"strings"
)

type listview struct {
	*tweetview
	list anaconda.List

	loading             lock
	loadNewTweetCh      chan []anaconda.Tweet
	loadIntervalTweetCh chan []anaconda.Tweet
}

func newListview() *listview {
	return &listview{
		tweetview:           newTweetview(),
		loadNewTweetCh:      make(chan []anaconda.Tweet),
		loadIntervalTweetCh: make(chan []anaconda.Tweet),
	}
}

func (lv *listview) loadTweet(sinceID int64) {
	if lv.loading.isLocking() {
		return
	}
	lv.loading.lock()
	defer lv.loading.unlock()
	changeBufferState("Loading List...")
	if err := lv.fetchList(); err != nil {
		changeBufferState("List Err")
		return
	}
	val := url.Values{}
	if sinceID > 0 {
		val.Add("since_id", strconv.FormatInt(sinceID, 10))
	}
	timeline, err := api.GetListTweets(lv.list.Id, true, val)
	if err != nil {
		changeBufferState("Err:Loading List")
		return
	}
	changeBufferState(fmt.Sprintf("Load!(%d tweets)", len(timeline)))
	lv.loadNewTweetCh <- timeline
}

func (lv *listview) loadIntervalTweet(maxID int64) {
	if lv.loading.isLocking() {
		return
	}
	lv.loading.lock()
	defer lv.loading.unlock()
	changeBufferState("Loading List...")
	if err := lv.fetchList(); err != nil {
		changeBufferState("List Err")
		return
	}
	val := url.Values{}
	if maxID > 0 {
		val.Add("max_id", strconv.FormatInt(maxID, 10))
	}
	timeline, err := api.GetListTweets(lv.list.Id, true, val)
	if err != nil {
		changeBufferState("Err:Loading List")
		return
	}
	changeBufferState(fmt.Sprintf("Load!(%d tweets)", len(timeline)))
	lv.loadIntervalTweetCh <- timeline
}

func (lv *listview) setListName(owner, name string) {
	lv.list = anaconda.List{Slug: name}
	lv.list.User.ScreenName = owner
	lv.tweets = []tweetstatus{tweetstatus{ReloadMark: true}}
}

func (lv *listview) setListID(id int64) {
	lv.list = anaconda.List{Id: id}
	lv.tweets = []tweetstatus{tweetstatus{ReloadMark: true}}
}

func (lv *listview) fetchList() error {
	if lv.list.Id == 0 || (lv.list.User.ScreenName == "" || lv.list.Slug == "") {
		val := url.Values{}
		if lv.list.Id != 0 {
			val.Add("list_id", strconv.FormatInt(lv.list.Id, 10))
		} else if lv.list.User.ScreenName != "" && lv.list.Slug != "" {
			val.Add("owner_screen_name", lv.list.User.ScreenName)
			val.Add("slug", lv.list.Slug)
		} else {
			return errors.New("can't fetch List information")
		}
		list, err := api.GetList(val)
		if err != nil {
			return err
		}
		lv.list = list
		return nil
	}
	return nil
}

func (lv *listview) draw() {
	tweets := lv.tweets
	if len(tweets) == 0 || tweets[0].ReloadMark {
		drawText("Now Loading...", 0, 0, ColorWhite, ColorBackground)
		return
	}

	width, _ := getTermSize()
	lines := strings.Split(runewidth.Wrap(lv.list.Description, width), "\n")

	lv.scrollOffset = 1 + len(lines)
	lv.tweetview.draw()

	text := lv.list.FullName
	x := 0
	y := 0
	drawText(text, x, y, ColorWhite, ColorGray2)
	x += runewidth.StringWidth(text)
	fillLine(x, y, ColorGray2)
	y++
	for _, t := range lines {
		drawText(t, 0, y, ColorWhite, ColorGray2)
		x = runewidth.StringWidth(t)
		fillLine(x, y, ColorGray2)
		y++
	}
}

func (lv *listview) resetScroll() {
	if len(lv.tweets) == 0 || lv.tweets[0].Content == nil {
		return
	}
	width, _ := getTermSize()
	lines := strings.Split(runewidth.Wrap(lv.list.Description, width), "\n")
	lv.scrollOffset = (1 + len(lines))

	lv.tweetview.resetScroll()
}
