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
	"strings"
	"time"
)

type view struct {
	timelineview     *timelineview
	mentionview      *mentionview
	conversationview *conversationview
	usertimelineview *usertimelineview
	favoriteview     *favoriteview
	listview         *listview
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
	view.favoriteview = newFavoriteview()
	view.listview = newListview()
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
	list
	favorite
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
	go func() {
		ticker := time.NewTicker(time.Second)
		counter := 0
		for {
			select {
			case <-ticker.C:
				if counter > 0 {
					counter--
					if counter == 0 {
						stateCh <- ""
					}
				}
			case v := <-stateClearCh:
				counter = v
			}
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
			view.conversationview.addPreviousTweet(wrapTweet(tw))
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
		case tw := <-view.favoriteview.loadNewTweetCh:
			tweetmap.registerTweets(tw)
			t := wrapTweets(tw)
			view.favoriteview.addNewTweet(t)
			view.refreshAll()
		case tw := <-view.favoriteview.loadIntervalTweetCh:
			tweetmap.registerTweets(tw)
			view.favoriteview.addIntervalTweet(wrapTweets(tw))
			view.refreshAll()
		case tw := <-view.listview.loadNewTweetCh:
			tweetmap.registerTweets(tw)
			t := wrapTweets(tw)
			view.listview.addNewTweet(t)
			view.refreshAll()
		case tw := <-view.listview.loadIntervalTweetCh:
			tweetmap.registerTweets(tw)
			view.listview.addIntervalTweet(wrapTweets(tw))
			view.refreshAll()
		case state := <-stateCh:
			if !view.buffer.inputing {
				view.buffer.setState(state)
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
	case favorite:
		view.buffer.linePosInfo = view.favoriteview.cursorPosition + 1
		view.favoriteview.draw()
	case list:
		view.buffer.linePosInfo = view.listview.cursorPosition + 1
		view.listview.draw()
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
	view.favoriteview.resetScroll()
	view.listview.resetScroll()
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
	case favorite:
		view.handleFavoriteMode(ev)
	case list:
		view.handleListMode(ev)
	}
}

func (view *view) handleCommonEvent(ev termbox.Event, tv *tweetview) {
	cursorPositionTweet := tv.tweets[tv.cursorPosition]
	switch view.handleAction(ev, KEYBIND_MODE_COMMON) {
	case ACTION_PREVIOUS_TWEET:
		tv.cursorDown()
	case ACTION_NEXT_TWEET:
		tv.cursorUp()
	case ACTION_TURN_INPUT_MODE:
		view.turnInputMode()
	case ACTION_LIKE_TWEET:
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
	case ACTION_RETWEET:
		if cursorPositionTweet.ReloadMark || cursorPositionTweet.Empty {
			return
		}
		if !cursorPositionTweet.isRetweeted() {
			cursorPositionTweet.setRetweeted(true)
			go retweet(cursorPositionTweet.Content.Id)
		}
	case ACTION_TURN_COMMAND_MODE:
		view.turnCommandMode()
	case ACTION_QUIT:
		view.quit = true
	case ACTION_TURN_CONVERSATION_VIEW_MODE:
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
	case ACTION_MENTION:
		if cursorPositionTweet.Empty || cursorPositionTweet.ReloadMark {
			return
		}
		view.turnReplyMode(cursorPositionTweet)
	case ACTION_TURN_USER_TIMELINE_MODE:
		if cursorPositionTweet.Empty || cursorPositionTweet.ReloadMark {
			return
		} else if view.usertimelineview.loading.isLocking() {
			return
		}
		t := cursorPositionTweet.Content
		if t.RetweetedStatus != nil {
			t = t.RetweetedStatus
		}
		view.turnUserTimelineMode(t.User.ScreenName)
	case ACTION_TURN_HOME_TIMELINE_MODE:
		view.turnHomeTimelineMode()
	case ACTION_TURN_MENTION_VIEW_MODE:
		view.turnMentionviewMode()
	case ACTION_OPEN_URL:
		if cursorPositionTweet.ReloadMark || cursorPositionTweet.Empty {
			return
		}
		for _, url := range cursorPositionTweet.Content.Entities.Urls {
			go openCommand(url.Expanded_url)
		}
	case ACTION_OPEN_IMAGES:
		if cursorPositionTweet.ReloadMark || cursorPositionTweet.Empty {
			return
		}
		urls := make([]string, 0)
		for _, media := range cursorPositionTweet.Content.ExtendedEntities.Media {
			urls = append(urls, media.Media_url_https)
		}
		if len(urls) > 0 {
			go openMedia(urls)
		}
	case ACTION_MOVE_TO_TOP_TWEET:
		tv.cursorMoveToTop()
	case ACTION_MOVE_TO_BOTTOM_TWEET:
		tv.cursorMoveToBottom()
	}
}

func (view *view) handleHometimelineMode(ev termbox.Event) {
	cursorPositionTweet := view.timelineview.
		tweets[view.timelineview.cursorPosition]

	switch view.handleAction(ev, KEYBIND_MODE_HOME_TIMELINE) { // go conversation view
	case ACTION_LOAD_PREVIOUSE_TWEETS:
		if cursorPositionTweet.ReloadMark {
			if !view.timelineview.isEmpty() {
				go view.timelineview.loadIntervalTweet(view.timelineview.tweets[view.
					timelineview.cursorPosition-1].Content.Id)
			} else {
				go view.timelineview.loadTweet(0)
			}
		}
	case ACTION_LOAD_NEW_TWEETS:
		if !view.timelineview.isEmpty() {
			go view.timelineview.loadTweet(view.timelineview.tweets[0].Content.Id)
		} else {
			go view.timelineview.loadTweet(0)
		}
	default:
		view.handleCommonEvent(ev, view.timelineview.tweetview)
	}
	view.refreshAll()
}

func (view *view) handleInputMode(ev termbox.Event) {
	switch view.handleAction(ev, KEYBIND_MODE_INPUT) {
	case ACTION_MOVE_LEFT:
		view.buffer.cursorMoveBackward()
	case ACTION_MOVE_RIGHT:
		view.buffer.cursorMoveForward()
	case ACTION_MOVE_UP:
		view.buffer.cursorMoveUp()
	case ACTION_MOVE_DOWN:
		view.buffer.cursorMoveDown()
	case ACTION_INSERT_SPACE:
		view.buffer.insertRune(' ')
	case ACTION_EXIT_INPUT_MODE:
		view.exitInputMode()
		view.refreshAll()
		return
	case ACTION_DELETE_RUNE:
		view.buffer.deleteRuneBackward()
	case ACTION_TEXT_CUT:
		view.buffer.cutToClipboard()
	case ACTION_TEXT_PASTE:
		view.buffer.pasteFromClipboard()
	case ACTION_MOVE_LINE_TOP:
		view.buffer.cursorMoveToLineTop()
	case ACTION_MOVE_LINE_BOTTOM:
		view.buffer.cursorMoveToLineBottom()
	case ACTION_TURN_CONFIRM_MODE:
		if len(view.buffer.content) != 0 {
			view.turnConfirmMode()
		}
	case ACTION_INSERT_NEW_LINE:
		if view.buffer.commanding {
			view.executeCommand(string(view.buffer.content))
			view.buffer.updateCursorPosition()
			view.refreshAll()
			return
		} else {
			view.buffer.insertLF()
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
	switch view.handleAction(ev, KEYBIND_MODE_CONFIRM) {
	case ACTION_CANCEL_SUBMIT:
		view.buffer.inputing = true
		view.buffer.confirm = false
		view.buffer.cursorMoveToLineBottom()
		view.buffer.updateCursorPosition()
	case ACTION_SUBMIT_TWEET:
		go view.buffer.process(string(view.buffer.content))
		view.exitConfirmMode()
	}
	view.refreshAll()
}

func (view *view) handleMentionviewMode(ev termbox.Event) {
	switch view.handleAction(ev, KEYBIND_MODE_MENTION_VIEW) {
	case ACTION_LOAD_PREVIOUSE_MENTIONS:
		if !view.mentionview.isEmpty() {
			go view.mentionview.loadIntervalTweet(view.mentionview.
				tweets[view.mentionview.cursorPosition-1].Content.Id)
		} else {
			go view.mentionview.loadTweet(0)
		}
	case ACTION_LOAD_NEW_MENTIONS:
		if !view.mentionview.isEmpty() {
			go view.mentionview.loadTweet(view.mentionview.tweets[0].Content.Id)
		} else {
			go view.mentionview.loadTweet(0)
		}
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

func (view *view) executeCommand(input string) {
	view.exitInputMode()
	splited := strings.SplitN(input, " ", 2)
	noArg := len(splited) < 2
	cmd := splited[0]
	args := ""
	if !noArg {
		args = strings.TrimSuffix(strings.TrimPrefix(splited[1], " "), " ")
	}
	// replace args when it is special keywords
	if args == "this" || args == "." {
		vm := view.getCurrentViewMode()
		if vm == usertimeline && view.usertimelineview.screenName != "" {
			args = view.usertimelineview.screenName
		} else if vm == favorite && view.favoriteview.screenName != "" {
			args = view.favoriteview.screenName
		}
	}
	switch cmd {
	case "user":
		if noArg {
			changeBufferState("Err! user command needs argument")
			return
		}
		view.turnUserTimelineMode(args)
	case "list":
		if noArg {
			changeBufferState("Err! list command needs argument")
			return
		}
		resplited := strings.Split(args, "/")
		var un, ln string
		switch len(resplited) {
		case 0:
			return
		case 1:
			un = user.ScreenName
			ln = resplited[0]
		case 2:
			un = resplited[0]
			ln = resplited[1]
		}
		view.turnListModeWithName(un, ln)
	case "favorite", "fav":
		if noArg {
			changeBufferState("Err! fav command needs argument")
			return
		}
		view.turnFavoriteviewMode(args)
	case "follow":
		if noArg {
			changeBufferState("Err! follow command needs argument")
			return
		}
		go func() {
			u, err := api.FollowUser(args)
			if err != nil {
				changeBufferState("Err: Couldn't follow specified user")
				return
			}
			changeBufferState("Succeed! following @" + u.ScreenName)
		}()
	case "unfollow":
		if noArg {
			changeBufferState("Err! unfollow command needs argument")
			return
		}
		go func() {
			u, err := api.UnfollowUser(args)
			if err != nil {
				changeBufferState("Err: Couldn't unfollow specified user")
				return
			}
			changeBufferState("Succeed! unfollowing @" + u.ScreenName)
		}()
	case "set_footer":
		if noArg {
			changeBufferState("Err! set_footer command needs argument")
			return
		}
		view.buffer.footer = args
	case "unset_footer":
		view.buffer.footer = ""
	default:
		changeBufferState("Commnad Err")
	}
}

func (view *view) handleConversationMode(ev termbox.Event) {
	switch view.handleAction(ev, KEYBIND_MODE_CONVERSATION) {
	case ACTION_EXIT_CONVERSATION_MODE:
		view.exitConversationviewMode()
	default:
		view.handleCommonEvent(ev, view.conversationview.tweetview)
	}
	view.refreshAll()
}

func (view *view) handleUserTimelineMode(ev termbox.Event) {
	cursorPositionTweet := view.usertimelineview.
		tweets[view.usertimelineview.cursorPosition]
	switch view.handleAction(ev, KEYBIND_MODE_USER_TIMELINE) {
	case ACTION_LOAD_PREVIOUSE_USER_TWEETS:
		if cursorPositionTweet.ReloadMark && view.usertimelineview.cursorPosition >= 1 {
			go view.usertimelineview.loadIntervalTweet(view.usertimelineview.
				tweets[view.usertimelineview.cursorPosition-1].Content.Id)
		}
	case ACTION_LOAD_NEW_USER_TWEETS:
		if !view.usertimelineview.loading.isLocking() {
			if !view.usertimelineview.isEmpty() {
				go view.usertimelineview.loadTweet(view.
					usertimelineview.tweets[0].Content.Id)
			} else {
				go view.usertimelineview.loadTweet(0)
			}
		}
	default:
		view.handleCommonEvent(ev, view.usertimelineview.tweetview)
	}

	view.refreshAll()
}

func (view *view) handleFavoriteMode(ev termbox.Event) {
	cursorPositionTweet := view.favoriteview.
		tweets[view.favoriteview.cursorPosition]
	switch view.handleAction(ev, KEYBIND_MODE_USER_FAVORITE) {
	case ACTION_LOAD_PREVIOUSE_USER_TWEETS:
		if cursorPositionTweet.ReloadMark && view.favoriteview.cursorPosition >= 1 {
			go view.favoriteview.loadIntervalTweet(view.favoriteview.
				tweets[view.favoriteview.cursorPosition-1].Content.Id)
		}
	case ACTION_LOAD_NEW_USER_TWEETS:
		if !view.favoriteview.loading.isLocking() {
			if !view.favoriteview.isEmpty() {
				go view.favoriteview.loadTweet(view.
					favoriteview.tweets[0].Content.Id)
			} else {
				go view.favoriteview.loadTweet(0)
			}
		}
	default:
		view.handleCommonEvent(ev, view.favoriteview.tweetview)
	}

	view.refreshAll()
}

func (view *view) handleListMode(ev termbox.Event) {
	cursorPositionTweet := view.listview.
		tweets[view.listview.cursorPosition]
	switch view.handleAction(ev, KEYBIND_MODE_LIST_VIEW) {
	case ACTION_LOAD_PREVIOUSE_LIST:
		if cursorPositionTweet.ReloadMark && view.listview.cursorPosition >= 1 {
			go view.listview.loadIntervalTweet(view.listview.
				tweets[view.listview.cursorPosition-1].Content.Id)
		}
	case ACTION_LOAD_NEW_LIST:
		go view.listview.loadTweet(view.
			listview.tweets[0].Content.Id)
	default:
		view.handleCommonEvent(ev, view.listview.tweetview)
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
	view.buffer.clear()
	view.setViewMode(mention)
	view.buffer.setModeStr(mention)
}

func (view *view) turnConversationviewMode() {
	view.buffer.clear()
	view.setViewMode(conversation)
	view.buffer.setModeStr(conversation)
	view.conversationview.cursorPosition = 0
	view.conversationview.scroll = 0
	go view.conversationview.loadTweet()
}

func (view *view) turnUserTimelineMode(screenName string) {
	view.usertimelineview.setUserScreenName(screenName)
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
func (view *view) turnFavoriteviewMode(screenName string) {
	view.favoriteview.setUserScreenName(screenName)
	view.setViewMode(favorite)
	view.buffer.setModeStr(favorite)
	view.favoriteview.cursorPosition = 0
	view.favoriteview.scroll = 0
	if view.favoriteview.isEmpty() {
		go view.favoriteview.loadTweet(0)
	} else if view.favoriteview.tweets[0].Content != nil {
		go view.
			favoriteview.loadTweet(view.
			favoriteview.tweets[0].Content.Id)
	}

}
func (view *view) turnListModeWithName(owner, name string) {
	view.listview.setListName(owner, name)
	view.setViewMode(list)
	view.buffer.setModeStr(list)
	view.listview.cursorPosition = 0
	view.listview.scroll = 0
	if view.listview.isEmpty() {
		go view.listview.loadTweet(0)
	}

}

func (view *view) turnInputMode() {
	view.buffer.inputing = true
	view.buffer.clear()
	// Set footer
	if view.buffer.footer != "" {
		// insert a SPACE before footer
		view.buffer.setContent(" " + view.buffer.footer)
	}
	view.buffer.cursorMoveToLineTop()
	view.buffer.updateCursorPosition()
	view.buffer.process = view.sendNewTweet
}

func (view *view) turnReplyMode(ts tweetstatus) {
	view.buffer.inputing = true
	view.buffer.setContent("@" + ts.Content.User.ScreenName + " ")
	view.buffer.cursorMoveToLineBottom()
	view.buffer.updateCursorPosition()
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

func (view *view) turnCommandMode() {
	view.buffer.inputing = true
	view.buffer.commanding = true
	view.buffer.clear()
	view.buffer.cursorMoveToLineBottom()
}

func (view *view) exitInputMode() {
	view.buffer.inputing = false
	view.buffer.commanding = false
	view.buffer.process = nil
	view.buffer.clear()
	view.buffer.setModeStr(view.getCurrentViewMode())
	view.buffer.cursorX = 0
	termbox.HideCursor()
}

func (view *view) exitConfirmMode() {
	view.buffer.inputing = false
	view.buffer.confirm = false
	view.buffer.process = nil
	view.buffer.clear()
	view.buffer.setModeStr(view.getCurrentViewMode())
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
	case list:
		view.setViewMode(list)
		view.buffer.setModeStr(list)
	case favorite:
		view.setViewMode(favorite)
		view.buffer.setModeStr(favorite)
	}
}
