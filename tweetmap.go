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
	"github.com/ChimeraCoder/anaconda"
	"sync"
)

// TweetMap caches loaded Tweet
type TweetMap struct {
	content map[int64]*anaconda.Tweet
	mutex   *sync.RWMutex
}

func newTweetMap() *TweetMap {
	return &TweetMap{
		content: make(map[int64]*anaconda.Tweet, 300),
		mutex:   new(sync.RWMutex),
	}
}

func (tm *TweetMap) registerTweet(tweet *anaconda.Tweet) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.content[tweet.Id] = tweet
}

func (tm *TweetMap) registerTweets(tweets []anaconda.Tweet) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	for i := range tweets {
		tm.content[tweets[i].Id] = &tweets[i]
	}
}

func (tm *TweetMap) get(id int64) (*anaconda.Tweet, bool) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	val, ok := tm.content[id]
	return val, ok
}
