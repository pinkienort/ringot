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

type favoriteview struct {
	*usertimelineview
}

func newFavoriteview() *favoriteview {
	return &favoriteview{
		usertimelineview: newUsertimelineview(),
	}
}

func (fv *favoriteview) loadTweet(sinceID int64) {
	if fv.loading.isLocking() {
		return
	}
	fv.loading.lock()
	defer fv.loading.unlock()
	changeBufferState("Loading...")

	if _, ok := profilemap.get(fv.screenName); !ok {
		u, err := api.GetUsersShow(fv.screenName, nil)
		if err == nil {
			profilemap.registerProfile(&u)
		} else {
			changeBufferState("Err:User profile Loading")
			return
		}
	}

	val := url.Values{}
	val.Add("count", strconv.Itoa(20))
	val.Add("screen_name", fv.screenName)
	if sinceID > 0 {
		val.Add("since_id", strconv.FormatInt(sinceID, 10))
	}
	timeline, err := api.GetFavorites(val)
	if err != nil {
		changeBufferState("Err:Loading")
		return
	}
	changeBufferState(fmt.Sprintf("Load!(%d tweets)", len(timeline)))
	fv.loadNewTweetCh <- timeline
}

func (fv *favoriteview) loadIntervalTweet(maxID int64) {
	if fv.loading.isLocking() {
		return
	}
	fv.loading.lock()
	defer fv.loading.unlock()
	changeBufferState("Loading...")
	val := url.Values{}
	val.Add("count", strconv.Itoa(CountTweet))
	val.Add("screen_name", fv.screenName)
	if maxID > 0 {
		val.Add("max_id", strconv.FormatInt(maxID, 10))
	}
	timeline, err := api.GetFavorites(val)
	if err != nil {
		changeBufferState("Err:Loading")
		return
	}
	changeBufferState(fmt.Sprintf("Load!(%d tweets)", len(timeline)))
	fv.loadIntervalTweetCh <- timeline
}

func (fv *favoriteview) addNewTweet(tss []tweetstatus) {
	if fv.userProfile == nil {
		u, ok := profilemap.get(fv.screenName)
		if ok {
			fv.userProfile = u
		}
	}
	if len(tss) == 0 {
		return
	}
	fv.tweetview.addNewTweet(tss)
}

func (fv *favoriteview) addIntervalTweet(tss []tweetstatus) {
	if fv.userProfile == nil {
		u, ok := profilemap.get(fv.screenName)
		if ok {
			fv.userProfile = u
		}
	}
	if len(tss) == 0 {
		return
	}
	fv.tweetview.addIntervalTweet(tss)
}
