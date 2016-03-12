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
)

const (
	loadMax = 50
)

type conversationview struct {
	*tweetview

	loadPreviousTweetCh chan *anaconda.Tweet
}

func newConversationview() *conversationview {
	return &conversationview{
		tweetview:           newTweetview(),
		loadPreviousTweetCh: make(chan *anaconda.Tweet),
	}
}

func (cv *conversationview) setTopTweet(ts tweetstatus) {
	if ts.Empty || ts.Content == nil {
		return
	}
	cv.tweets = []tweetstatus{ts}

	count := 1
	tweet := ts.Content
	for {
		id := tweet.InReplyToStatusID
		if id == 0 {
			break
		} else if t, ok := tweetmap.get(id); ok {
			cv.tweets = append(cv.tweets, tweetstatus{Content: t})
			tweet = t
		} else {
			break
		}
		count++
		if count >= loadMax {
			break
		}
	}
}

func (cv *conversationview) addPreviousTweet(ts tweetstatus) {
	cv.tweets = append(cv.tweets, ts)
}

func (cv *conversationview) loadTweet() {
	count := 1
	if len(cv.tweets) < 1 {
		return
	}
	tweet := cv.tweets[len(cv.tweets)-1].Content
	for {
		id := tweet.InReplyToStatusID
		if id == 0 {
			break
		} else if t, ok := tweetmap.get(id); ok {
			cv.loadPreviousTweetCh <- t
			tweet = t
		} else {
			t, err := api.GetTweet(id, nil)
			if err != nil {
				changeBufferState(fmt.Sprintf("Err:Load Tweet(ID:%d)", id))
			}
			cv.loadPreviousTweetCh <- &t
			tweet = &t

		}
		count++
		if count >= loadMax {
			break
		}
	}
}
