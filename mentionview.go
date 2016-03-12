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
	"net/url"
	"strconv"
)

type mentionview struct {
	*timelineview
}

func newMentionview() *mentionview {
	return &mentionview{
		timelineview: newTimelineview(),
	}
}

func (mv *mentionview) loadTweet(sinceID int64) {
	if mv.loading.isLocking() {
		return
	}
	mv.loading.lock()
	defer mv.loading.unlock()
	changeBufferState("Mention Loading...")
	val := url.Values{}
	val.Add("count", strconv.Itoa(CountTweet))
	if sinceID > 0 {
		val.Add("since_id", strconv.FormatInt(sinceID, 10))
	}
	timeline, err := api.GetMentionsTimeline(val)
	if err != nil {
		changeBufferState("Err:Loading")
		return
	}
	changeBufferState(fmt.Sprintf("Load!(%d tweets)", len(timeline)))
	mv.loadNewTweetCh <- timeline
}

func (mv *mentionview) loadIntervalTweet(maxID int64) {
	if mv.loading.isLocking() {
		return
	}
	mv.loading.lock()
	defer mv.loading.unlock()
	changeBufferState("Loading...")
	val := url.Values{}
	val.Add("count", strconv.Itoa(CountTweet))
	if maxID > 0 {
		val.Add("max_id", strconv.FormatInt(maxID, 10))
	}
	timeline, err := api.GetMentionsTimeline(val)
	if err != nil {
		changeBufferState("Err:Loading")
		return
	}
	changeBufferState(fmt.Sprintf("Load!(%d tweets)", len(timeline)))
	mv.loadIntervalTweetCh <- timeline
}
