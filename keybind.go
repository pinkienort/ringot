// This file is part of Ringot.
/*
Copyright 2016 pinkienort <cantabilehisa@gmail.com>

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

import termbox "github.com/nsf/termbox-go"

type KeybindMode uint8

const (
	KEYBIND_MODE_COMMON = iota
	KEYBIND_MODE_CONVERSATION
	KEYBIND_MODE_HOME_TIMELINE
	KEYBIND_MODE_INPUT
	KEYBIND_MODE_CONFIRM
	KEYBIND_MODE_MENTION_VIEW
	KEYBIND_MODE_USER_TIMELINE
	KEYBIND_MODE_USER_FAVORITE
	KEYBIND_MODE_LIST_VIEW
)

type Action uint8
type keybind struct {
	Mod    termbox.Modifier
	Key    termbox.Key
	Ch     rune
	Action Action
}

const ( /* common event action list */
	ACTION_LIKE_TWEET = iota + 1
	ACTION_MENTION
	ACTION_RETWEET
	ACTION_OPEN_IMAGES
	ACTION_OPEN_USER_PROFILE_IMAGE
	ACTION_NEXT_TWEET
	ACTION_PREVIOUS_TWEET
	ACTION_PAGE_DOWN
	ACTION_PAGE_UP
	ACTION_MOVE_TO_TOP_TWEET
	ACTION_MOVE_TO_BOTTOM_TWEET
	ACTION_TURN_INPUT_MODE
	ACTION_TURN_COMMAND_MODE
	ACTION_TURN_HOME_TIMELINE_MODE
	ACTION_TURN_CONVERSATION_VIEW_MODE
	ACTION_TURN_MENTION_VIEW_MODE
	ACTION_TURN_USER_TIMELINE_MODE
	ACTION_QUIT
	ACTION_OPEN_URL
	ACTION_SHOW_HELP
)
const ( /* home timeline action list */
	ACTION_LOAD_PREVIOUSE_TWEETS = iota + 1
	ACTION_LOAD_NEW_TWEETS
)
const ( /* input mode action list */
	ACTION_MOVE_LEFT = iota + 1
	ACTION_MOVE_RIGHT
	ACTION_MOVE_UP
	ACTION_MOVE_DOWN
	ACTION_INSERT_SPACE
	ACTION_EXIT_INPUT_MODE
	ACTION_DELETE_RUNE
	ACTION_MOVE_LINE_TOP
	ACTION_MOVE_LINE_BOTTOM
	ACTION_TURN_CONFIRM_MODE
	ACTION_INSERT_NEW_LINE
	ACTION_TEXT_CUT
	ACTION_TEXT_PASTE
)
const ( /* confirm mode action list */
	ACTION_CANCEL_SUBMIT = iota + 1
	ACTION_SUBMIT_TWEET
)
const ( /* mention view mode action list */
	ACTION_LOAD_PREVIOUSE_MENTIONS = iota + 1
	ACTION_LOAD_NEW_MENTIONS
)
const ( /* conversation mode action list */
	ACTION_EXIT_CONVERSATION_MODE = iota + 1
)
const ( /* user timeline mode action list */
	ACTION_LOAD_PREVIOUSE_USER_TWEETS = iota + 1
	ACTION_LOAD_NEW_USER_TWEETS
)
const ( /* list mode actin list */
	ACTION_LOAD_PREVIOUSE_LIST = iota + 1
	ACTION_LOAD_NEW_LIST
)

const NO_MOD = 0
const NO_KEY = 0
const NO_CH = 0
const NO_ACTION = 0

/* keybind list */
var commonKeybindList = []keybind{
	{NO_MOD, termbox.KeyCtrlS, NO_CH, ACTION_TURN_INPUT_MODE},
	{NO_MOD, termbox.KeyCtrlW, NO_CH, ACTION_MENTION},
	{NO_MOD, termbox.KeyCtrlF, NO_CH, ACTION_LIKE_TWEET},
	{NO_MOD, termbox.KeyCtrlV, NO_CH, ACTION_RETWEET},
	{NO_MOD, termbox.KeyCtrlO, NO_CH, ACTION_OPEN_URL},
	{NO_MOD, termbox.KeyCtrlP, NO_CH, ACTION_OPEN_IMAGES},

	{NO_MOD, termbox.KeyCtrlZ, NO_CH, ACTION_TURN_HOME_TIMELINE_MODE},
	{NO_MOD, termbox.KeyCtrlX, NO_CH, ACTION_TURN_MENTION_VIEW_MODE},
	{NO_MOD, termbox.KeyCtrlD, NO_CH, ACTION_TURN_USER_TIMELINE_MODE},
	{NO_MOD, termbox.KeyArrowRight, NO_CH, ACTION_TURN_CONVERSATION_VIEW_MODE},
	{termbox.ModAlt, NO_KEY, 'x', ACTION_TURN_COMMAND_MODE}, /* TODO: need ModAlt field */
	{NO_MOD, termbox.KeyCtrlQ, NO_CH, ACTION_QUIT},

	{NO_MOD, termbox.KeyArrowUp, NO_CH, ACTION_NEXT_TWEET},
	{NO_MOD, termbox.KeyArrowDown, NO_CH, ACTION_PREVIOUS_TWEET},
	{NO_MOD, termbox.KeyHome, NO_CH, ACTION_MOVE_TO_TOP_TWEET},
	{NO_MOD, termbox.KeyEnd, NO_CH, ACTION_MOVE_TO_BOTTOM_TWEET},
	{NO_MOD, termbox.KeyPgup, NO_CH, ACTION_PAGE_UP},
	{NO_MOD, termbox.KeyPgdn, NO_CH, ACTION_PAGE_DOWN},
}

var homeTimelineKeybindList = []keybind{
	{NO_MOD, termbox.KeyEnter, NO_CH, ACTION_LOAD_PREVIOUSE_MENTIONS},
	{NO_MOD, termbox.KeySpace, NO_CH, ACTION_LOAD_PREVIOUSE_TWEETS},
	{NO_MOD, termbox.KeyCtrlR, NO_CH, ACTION_LOAD_NEW_TWEETS},
}

var inputModeKeybindList = []keybind{
	{NO_MOD, termbox.KeyArrowLeft, NO_CH, ACTION_MOVE_LEFT},
	{NO_MOD, termbox.KeyArrowRight, NO_CH, ACTION_MOVE_RIGHT},
	{NO_MOD, termbox.KeyArrowUp, NO_CH, ACTION_MOVE_UP},
	{NO_MOD, termbox.KeyArrowDown, NO_CH, ACTION_MOVE_DOWN},
	{NO_MOD, termbox.KeySpace, NO_CH, ACTION_INSERT_SPACE},
	{NO_MOD, termbox.KeyEsc, NO_CH, ACTION_EXIT_INPUT_MODE},
	{NO_MOD, termbox.KeyCtrlG, NO_CH, ACTION_EXIT_INPUT_MODE},
	{NO_MOD, termbox.KeyBackspace, NO_CH, ACTION_DELETE_RUNE},
	{NO_MOD, termbox.KeyBackspace2, NO_CH, ACTION_DELETE_RUNE},
	{NO_MOD, termbox.KeyCtrlA, NO_CH, ACTION_MOVE_LINE_TOP},
	{NO_MOD, termbox.KeyCtrlE, NO_CH, ACTION_MOVE_LINE_BOTTOM},
	{NO_MOD, termbox.KeyCtrlJ, NO_CH, ACTION_TURN_CONFIRM_MODE},
	{NO_MOD, termbox.KeyEnter, NO_CH, ACTION_INSERT_NEW_LINE},
	{NO_MOD, termbox.KeyCtrlW, NO_CH, ACTION_TEXT_CUT},
	{NO_MOD, termbox.KeyCtrlY, NO_CH, ACTION_TEXT_PASTE},
}

var confirmModeKeybindList = []keybind{
	{NO_MOD, termbox.KeyEsc, NO_CH, ACTION_CANCEL_SUBMIT},
	{NO_MOD, termbox.KeyCtrlG, NO_CH, ACTION_CANCEL_SUBMIT},
	{NO_MOD, termbox.KeyEnter, NO_CH, ACTION_SUBMIT_TWEET},
}

var mentionViewModeKeybindList = []keybind{
	{NO_MOD, termbox.KeyEnter, NO_CH, ACTION_LOAD_PREVIOUSE_MENTIONS},
	{NO_MOD, termbox.KeySpace, NO_CH, ACTION_LOAD_PREVIOUSE_MENTIONS},
	{NO_MOD, termbox.KeyCtrlR, NO_CH, ACTION_LOAD_NEW_MENTIONS},
}

var conversationModeKeybindList = []keybind{
	{NO_MOD, termbox.KeyArrowLeft, NO_CH, ACTION_EXIT_CONVERSATION_MODE},
}

var userTimelineModeKeybindList = []keybind{
	{NO_MOD, termbox.KeyEnter, NO_CH, ACTION_LOAD_PREVIOUSE_USER_TWEETS},
	{NO_MOD, termbox.KeySpace, NO_CH, ACTION_LOAD_PREVIOUSE_USER_TWEETS},
	{NO_MOD, termbox.KeyCtrlR, NO_CH, ACTION_LOAD_NEW_USER_TWEETS},
	{NO_MOD, termbox.KeyCtrl8, NO_CH, ACTION_OPEN_USER_PROFILE_IMAGE},
}

var favoriteModeKeybindList = []keybind{
	{NO_MOD, termbox.KeyEnter, NO_CH, ACTION_LOAD_PREVIOUSE_USER_TWEETS},
	{NO_MOD, termbox.KeySpace, NO_CH, ACTION_LOAD_PREVIOUSE_USER_TWEETS},
	{NO_MOD, termbox.KeyCtrlR, NO_CH, ACTION_LOAD_NEW_USER_TWEETS},
}

var listModeKeybindList = []keybind{
	{NO_MOD, termbox.KeyEnter, NO_CH, ACTION_LOAD_PREVIOUSE_LIST},
	{NO_MOD, termbox.KeySpace, NO_CH, ACTION_LOAD_PREVIOUSE_LIST},
	{NO_MOD, termbox.KeyCtrlR, NO_CH, ACTION_LOAD_NEW_LIST},
}

func (view *view) handleAction(ev termbox.Event, mode KeybindMode) Action {
	var keybindList []keybind
	switch mode {
	case KEYBIND_MODE_COMMON:
		keybindList = commonKeybindList
	case KEYBIND_MODE_CONVERSATION:
		keybindList = conversationModeKeybindList
	case KEYBIND_MODE_HOME_TIMELINE:
		keybindList = homeTimelineKeybindList
	case KEYBIND_MODE_INPUT:
		keybindList = inputModeKeybindList
	case KEYBIND_MODE_CONFIRM:
		keybindList = confirmModeKeybindList
	case KEYBIND_MODE_MENTION_VIEW:
		keybindList = mentionViewModeKeybindList
	case KEYBIND_MODE_USER_TIMELINE:
		keybindList = userTimelineModeKeybindList
	case KEYBIND_MODE_USER_FAVORITE:
		keybindList = favoriteModeKeybindList
	case KEYBIND_MODE_LIST_VIEW:
		keybindList = listModeKeybindList
	}
	for i := 0; i < len(keybindList); i++ {
		if (ev.Mod == keybindList[i].Mod) &&
			(ev.Key == keybindList[i].Key) &&
			(ev.Ch == keybindList[i].Ch) {
			return keybindList[i].Action
		}
	}
	return NO_ACTION
}
