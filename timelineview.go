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
	"net/url"
	"strconv"
)

type timelineview struct {
	*tweetview

	loading             lock
	loadNewTweetCh      chan []anaconda.Tweet
	loadIntervalTweetCh chan []anaconda.Tweet
}

func newTimelineview() *timelineview {
	return &timelineview{
		tweetview:           newTweetview(),
		loadNewTweetCh:      make(chan []anaconda.Tweet),
		loadIntervalTweetCh: make(chan []anaconda.Tweet),
	}
}

func (tv *timelineview) loadTweet(sinceID int64) {
	if tv.loading.isLocking() {
		return
	}
	tv.loading.lock()
	defer tv.loading.unlock()
	changeBufferState("Loading...")
	val := url.Values{}
	val.Add("count", strconv.Itoa(CountTweet))
	if sinceID > 0 {
		val.Add("since_id", strconv.FormatInt(sinceID, 10))
	}
	timeline, err := api.GetHomeTimeline(val)
	if err != nil {
		changeBufferState("Err:Loading")
		return
	}
	changeBufferState(fmt.Sprintf("Load!(%d tweets)", len(timeline)))
	tv.loadNewTweetCh <- timeline
}

func (tv *timelineview) loadIntervalTweet(maxID int64) {
	if tv.loading.isLocking() {
		return
	}
	tv.loading.lock()
	defer tv.loading.unlock()
	changeBufferState("Loading...")
	val := url.Values{}
	val.Add("count", strconv.Itoa(50))
	if maxID > 0 {
		val.Add("max_id", strconv.FormatInt(maxID, 10))
	}
	timeline, err := api.GetHomeTimeline(val)
	if err != nil {
		changeBufferState("Err:Loading")
		return
	}
	changeBufferState(fmt.Sprintf("Load!(%d tweets)", len(timeline)))
	tv.loadIntervalTweetCh <- timeline
}
