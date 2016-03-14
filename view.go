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
	"github.com/nsf/termbox-go"
	"net/url"
	"strconv"
	"time"
)

type view struct {
	timelineview     *timelineview
	mentionview      *mentionview
	conversationview *conversationview
	usertimelineview *usertimelineview
	buffer           *buffer

	modeHistory []viewmode
	quit        bool
}

func newView() *view {
	view := &view{
		modeHistory: []viewmode{home},
		quit:        false,
	}
	view.timelineview = newTimelineview()
	view.mentionview = newMentionview()
	view.conversationview = newConversationview()
	view.usertimelineview = newUsertimelineview()
	view.buffer = newBuffer()
	return view
}

type (
	viewmode int
)

const (
	home viewmode = iota
	usertimeline
	mention
	conversation
)

func (view *view) Init() {
	view.initHomeTimeline()
	view.initMention()
	view.turnHomeTimelineMode()
	view.refreshAll()
}

func (view *view) initHomeTimeline() {
	ht, err := api.GetHomeTimeline(nil)
	if err != nil {
		changeBufferState("Can't initialize Home Timeline")
		return
	}
	tweetmap.registerTweets(ht)
	t := wrapTweets(ht)
	view.timelineview.addNewTweet(t)
}

func (view *view) initMention() {
	mt, err := api.GetMentionsTimeline(nil)
	if err != nil {
		return
	}
	tweetmap.registerTweets(mt)
	t := wrapTweets(mt)
	view.mentionview.addNewTweet(t)
}

func (view *view) Loop() {
	evCh := make(chan termbox.Event)
	go func() {
		for {
			evCh <- termbox.PollEvent()
		}
	}()
	for {
		select {
		case ev := <-evCh:
			if ev.Type == termbox.EventResize {
				setTermSize(ev.Width, ev.Height)
				view.resetScrollAll()
				view.buffer.updateCursorPosition()
				view.refreshAll()
			} else {
				view.handleEvent(ev)
			}
		case tw := <-view.timelineview.loadNewTweetCh:
			tweetmap.registerTweets(tw)
			view.timelineview.addNewTweet(wrapTweets(tw))
			view.refreshAll()
		case tw := <-view.timelineview.loadIntervalTweetCh:
			tweetmap.registerTweets(tw)
			view.timelineview.addIntervalTweet(wrapTweets(tw))
			view.refreshAll()
		case tw := <-view.mentionview.loadNewTweetCh:
			tweetmap.registerTweets(tw)
			view.mentionview.addNewTweet(wrapTweets(tw))
			view.refreshAll()
		case tw := <-view.mentionview.loadIntervalTweetCh:
			tweetmap.registerTweets(tw)
			view.mentionview.addIntervalTweet(wrapTweets(tw))
			view.refreshAll()
		case tw := <-view.conversationview.loadPreviousTweetCh:
			tweetmap.registerTweet(tw)
			view.conversationview.addPreviousTweet(tweetstatus{Content: tw})
			view.refreshAll()
		case tw := <-view.usertimelineview.loadNewTweetCh:
			tweetmap.registerTweets(tw)
			t := wrapTweets(tw)
			view.usertimelineview.addNewTweet(t)
			view.refreshAll()
		case tw := <-view.usertimelineview.loadIntervalTweetCh:
			tweetmap.registerTweets(tw)
			view.usertimelineview.addIntervalTweet(wrapTweets(tw))
			view.refreshAll()
		case state := <-stateCh:
			if !view.buffer.inputing {
				view.buffer.setContent(state)
				view.refreshBuffer()
			}
		}
		if view.quit {
			break
		}
	}

}

func (view *view) refreshAll() {
	termbox.Clear(ColorBackground, ColorBackground)

	switch view.getCurrentViewMode() {
	case home:
		view.buffer.linePosInfo = view.timelineview.cursorPosition + 1
		view.timelineview.draw()
	case mention:
		view.buffer.linePosInfo = view.mentionview.cursorPosition + 1
		view.mentionview.draw()
	case conversation:
		view.buffer.linePosInfo = view.conversationview.cursorPosition + 1
		view.conversationview.draw()
	case usertimeline:
		view.buffer.linePosInfo = view.usertimelineview.cursorPosition + 1
		view.usertimelineview.draw()
	}
	view.buffer.draw()
	termbox.Flush()
}

func (view *view) refreshBuffer() {
	view.buffer.draw()
	termbox.Flush()
}

func (view *view) resetScrollAll() {
	view.timelineview.resetScroll()
	view.mentionview.resetScroll()
	view.conversationview.resetScroll()
	view.usertimelineview.resetScroll()
}

func (view *view) handleEvent(ev termbox.Event) {
	if view.buffer.inputing {
		if view.buffer.confirm {
			view.handleConfirmMode(ev)
		} else {
			view.handleInputMode(ev)
		}
		return
	}
	switch view.getCurrentViewMode() {
	case home:
		view.handleHometimelineMode(ev)
	case mention:
		view.handleMentionviewMode(ev)
	case conversation:
		view.handleConversationMode(ev)
	case usertimeline:
		view.handleUserTimelineMode(ev)
	}
}

func (view *view) handleCommonEvent(ev termbox.Event, tv *tweetview) {
	cursorPositionTweet := tv.tweets[tv.cursorPosition]
	switch ev.Key {
	case termbox.KeyCtrlQ:
		view.quit = true
	case termbox.KeyArrowUp:
		tv.cursorUp()
	case termbox.KeyArrowDown:
		tv.cursorDown()
	case termbox.KeyArrowRight:
		if cursorPositionTweet.Empty || cursorPositionTweet.ReloadMark ||
			cursorPositionTweet.Content == nil {
			return
		}
		t := cursorPositionTweet.Content
		if t.RetweetedStatus != nil {
			t = t.RetweetedStatus
		}
		if t.InReplyToStatusID == 0 {
			return
		}
		view.conversationview.setTopTweet(tweetstatus{Content: t})
		view.turnConversationviewMode()
	case termbox.KeyCtrlS:
		view.turnTweetMode()
	case termbox.KeyCtrlW:
		if cursorPositionTweet.Empty || cursorPositionTweet.ReloadMark {
			return
		}
		view.turnReplyMode(cursorPositionTweet)
	case termbox.KeyCtrlD:
		if cursorPositionTweet.Empty || cursorPositionTweet.ReloadMark {
			return
		} else if view.usertimelineview.loading.isLocking() {
			return
		}
		t := cursorPositionTweet.Content
		if t.RetweetedStatus != nil {
			t = t.RetweetedStatus
		}
		view.usertimelineview.setUserScreenName(t.User.ScreenName)
		view.turnUserTimelineMode()
	case termbox.KeyCtrlZ:
		view.turnHomeTimelineMode()
	case termbox.KeyCtrlX:
		view.turnMentionviewMode()
	case termbox.KeyCtrlF:
		if cursorPositionTweet.ReloadMark || cursorPositionTweet.Empty {
			return
		}
		if !cursorPositionTweet.isFavorited() {
			cursorPositionTweet.setFavorited(true)
			go favoriteTweet(cursorPositionTweet.Content.Id)
		} else {
			cursorPositionTweet.setFavorited(false)
			go unfavoriteTweet(cursorPositionTweet.Content.Id)
		}
	case termbox.KeyCtrlV:
		if cursorPositionTweet.ReloadMark || cursorPositionTweet.Empty {
			return
		}
		if !cursorPositionTweet.isRetweeted() {
			cursorPositionTweet.setRetweeted(true)
			go retweet(cursorPositionTweet.Content.Id)
		}
	case termbox.KeyCtrlO:
		if cursorPositionTweet.ReloadMark || cursorPositionTweet.Empty {
			return
		}
		for _, url := range cursorPositionTweet.Content.Entities.Urls {
			go openCommand(url.Expanded_url)
		}
	case termbox.KeyCtrlP:
		if cursorPositionTweet.ReloadMark || cursorPositionTweet.Empty {
			return
		}
		for _, media := range cursorPositionTweet.Content.ExtendedEntities.Media {
			go openMedia(media.Media_url_https)
		}
	}
}

func (view *view) handleHometimelineMode(ev termbox.Event) {
	cursorPositionTweet := view.timelineview.
		tweets[view.timelineview.cursorPosition]
	switch ev.Key {
	case termbox.KeyEnter, termbox.KeySpace:
		if cursorPositionTweet.ReloadMark {
			if !view.timelineview.isEmpty() {
				go view.timelineview.loadIntervalTweet(view.timelineview.tweets[view.
					timelineview.cursorPosition-1].Content.Id)
			} else {
				go view.timelineview.loadTweet(0)
			}

		}
	case termbox.KeyCtrlR:
		if !view.timelineview.isEmpty() {
			go view.timelineview.loadTweet(view.timelineview.tweets[0].Content.Id)
		} else {
			go view.timelineview.loadTweet(0)
		}
	case termbox.KeyCtrlZ:
		// Do nothing
	default:
		view.handleCommonEvent(ev, view.timelineview.tweetview)
	}
	view.refreshAll()
}

func (view *view) handleInputMode(ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyArrowLeft:
		view.buffer.cursorMoveBackward()
	case termbox.KeyArrowRight:
		view.buffer.cursorMoveForward()
	case termbox.KeySpace:
		view.buffer.insertRune(' ')
	case termbox.KeyEsc, termbox.KeyCtrlG:
		view.exitTweetMode()
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		view.buffer.deleteRuneBackward()
	case termbox.KeyCtrlA:
		view.buffer.cursorMoveToTop()
	case termbox.KeyCtrlE:
		view.buffer.cursorMoveToBottom()
	case termbox.KeyCtrlJ:
		if len(view.buffer.content) != 0 {
			view.turnConfirmMode()
		}
	default:
		if ev.Ch != 0 {
			view.buffer.insertRune(ev.Ch)
		}
	}
	view.buffer.updateCursorPosition()
	view.refreshBuffer()
}

func (view *view) handleConfirmMode(ev termbox.Event) {
	if view.buffer.confirmLock.isLocking() {
		return
	}
	switch ev.Key {
	case termbox.KeyEsc, termbox.KeyCtrlG:
		view.buffer.inputing = true
		view.buffer.confirm = false
		view.buffer.cursorMoveToBottom()
		view.buffer.updateCursorPosition()
	case termbox.KeyEnter:
		go view.buffer.process(string(view.buffer.content))
		view.exitConfirmMode()
	}
	view.refreshBuffer()
}

func (view *view) handleMentionviewMode(ev termbox.Event) {
	cursorPositionTweet := view.mentionview.
		tweets[view.mentionview.cursorPosition]
	switch ev.Key {
	case termbox.KeyEnter, termbox.KeySpace:
		if cursorPositionTweet.ReloadMark {
			if !view.mentionview.isEmpty() {
				go view.mentionview.loadIntervalTweet(view.mentionview.
					tweets[view.mentionview.cursorPosition-1].Content.Id)
			} else {
				go view.mentionview.loadTweet(0)
			}
		}
	case termbox.KeyCtrlR:
		if !view.mentionview.isEmpty() {
			go view.mentionview.loadTweet(view.mentionview.tweets[0].Content.Id)
		} else {
			go view.mentionview.loadTweet(0)
		}
	case termbox.KeyCtrlX:
		// Do nothing
	default:
		view.handleCommonEvent(ev, view.mentionview.tweetview)
	}
	view.refreshAll()
}

func (view *view) sendNewTweet(status string) {
	if view.timelineview.loading.isLocking() || len(status) == 0 {
		return
	}
	view.timelineview.loading.lock()
	defer view.timelineview.loading.unlock()
	changeBufferState("Posting Tweet...")
	_, err := api.PostTweet(status, nil)
	if err != nil {
		changeBufferState("Err! Failed to tweet")
		return
	}
	changeBufferState("Tweet!")
}

func (view *view) handleConversationMode(ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyArrowLeft:
		view.exitConversationviewMode()
	case termbox.KeyArrowRight:
		// Do nothing
	default:
		view.handleCommonEvent(ev, view.conversationview.tweetview)
	}
	view.refreshAll()
}

func (view *view) handleUserTimelineMode(ev termbox.Event) {
	cursorPositionTweet := view.usertimelineview.
		tweets[view.usertimelineview.cursorPosition]
	switch ev.Key {
	case termbox.KeyEnter, termbox.KeySpace:
		if cursorPositionTweet.ReloadMark && view.usertimelineview.cursorPosition >= 1 {
			go view.usertimelineview.loadIntervalTweet(view.usertimelineview.
				tweets[view.usertimelineview.cursorPosition-1].Content.Id)
		}
	case termbox.KeyCtrlR:
		go view.usertimelineview.loadTweet(view.
			usertimelineview.tweets[0].Content.Id)
	default:
		view.handleCommonEvent(ev, view.usertimelineview.tweetview)
	}

	view.refreshAll()
}

func (view *view) setViewMode(mode viewmode) {
	view.modeHistory = append(view.modeHistory, mode)
	if len(view.modeHistory) > 5 {
		l := len(view.modeHistory)
		view.modeHistory = view.modeHistory[l-5 : l]
	}
}

func (view *view) getCurrentViewMode() viewmode {
	return view.getPreviousViewMode(0)
}

func (view *view) getPreviousViewMode(i int) viewmode {
	l := len(view.modeHistory)
	return view.modeHistory[l-i-1]
}

func (view *view) turnHomeTimelineMode() {
	view.setViewMode(home)
	view.buffer.setModeStr(home)
}

func (view *view) turnMentionviewMode() {
	view.buffer.setContent("")
	view.setViewMode(mention)
	view.buffer.setModeStr(mention)
}

func (view *view) turnConversationviewMode() {
	view.buffer.setContent("")
	view.setViewMode(conversation)
	view.buffer.setModeStr(conversation)
	view.conversationview.cursorPosition = 0
	view.conversationview.scroll = 0
	go view.conversationview.loadTweet()
}

func (view *view) turnUserTimelineMode() {
	view.setViewMode(usertimeline)
	view.buffer.setModeStr(usertimeline)
	view.usertimelineview.cursorPosition = 0
	view.usertimelineview.scroll = 0
	if view.usertimelineview.isEmpty() {
		go view.usertimelineview.loadTweet(0)
	} else if view.usertimelineview.tweets[0].Content != nil {
		go view.
			usertimelineview.loadTweet(view.
			usertimelineview.tweets[0].Content.Id)
	}

}

func (view *view) turnTweetMode() {
	view.buffer.inputing = true
	view.buffer.setContent("")
	view.buffer.cursorMoveToTop()
	view.buffer.updateCursorPosition()
	view.buffer.process = view.sendNewTweet
}

func (view *view) turnReplyMode(ts tweetstatus) {
	view.buffer.inputing = true
	view.buffer.setContent("@" + ts.Content.User.ScreenName + " ")
	view.buffer.cursorMoveToBottom()
	view.buffer.process = func(status string) {
		if view.timelineview.loading.isLocking() || len(status) == 0 {
			return
		}
		changeBufferState("Posting Tweet...")
		val := url.Values{}
		val.Add("in_reply_to_status_id", strconv.FormatInt(ts.Content.Id, 10))
		_, err := api.PostTweet(status, val)
		if err != nil {
			changeBufferState("Err! Failed to tweet")
			return
		}
		changeBufferState("Tweet!")
	}
}

func (view *view) turnConfirmMode() {
	termbox.HideCursor()
	view.buffer.confirm = true
	view.buffer.confirmLock.lock()
	go func() {
		// Wait a half second for confirm to tweet
		time.Sleep(time.Millisecond * 500)
		view.buffer.confirmLock.unlock()
	}()
}

func (view *view) exitTweetMode() {
	view.buffer.inputing = false
	view.buffer.process = nil
	view.buffer.setModeStr(view.getCurrentViewMode())
	view.buffer.setContent("")
	view.buffer.cursorX = 0
	termbox.HideCursor()
}

func (view *view) exitConfirmMode() {
	view.buffer.inputing = false
	view.buffer.confirm = false
	view.buffer.process = nil
	view.buffer.setModeStr(view.getCurrentViewMode())
	view.buffer.setContent("")
	view.buffer.cursorX = 0
}

func (view *view) exitConversationviewMode() {
	switch view.getPreviousViewMode(1) {
	case home:
		view.setViewMode(home)
		view.buffer.setModeStr(home)
	case mention:
		view.setViewMode(mention)
		view.buffer.setModeStr(mention)
	case usertimeline:
		view.setViewMode(usertimeline)
		view.buffer.setModeStr(usertimeline)
	}
}
