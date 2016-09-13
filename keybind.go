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
	KEYBIND_MODE_LIST_VIEW
)

type Action uint8

type keybind struct {
	Key    termbox.Key
	Ch     rune
	Action Action
}

const NO_ACTION = 0xff - 1

/* common event action list */
const (
	ACTION_LIKE_TWEET = iota
	ACTION_MENTION
	ACTION_RETWEET
	ACTION_OPEN_IMAGES
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
var commonKeybindList = []keybind {
	{ termbox.KeyCtrlS		,	0	,	ACTION_TURN_INPUT_MODE				},
	{ termbox.KeyCtrlW		,	0	,	ACTION_MENTION						},
	{ termbox.KeyCtrlF		,	0	,	ACTION_LIKE_TWEET					},
	{ termbox.KeyCtrlV		,	0	,	ACTION_RETWEET						},
	{ termbox.KeyCtrlO		,	0	,	ACTION_OPEN_URL						},
	{ termbox.KeyCtrlP		,	0	,	ACTION_OPEN_IMAGES					},

	{ termbox.KeyCtrlZ		,	0	,	ACTION_TURN_HOME_TIMELINE_MODE		},
	{ termbox.KeyCtrlX		,	0	,	ACTION_TURN_MENTION_VIEW_MODE		},
	{ termbox.KeyCtrlD		,	0	,	ACTION_TURN_USER_TIMELINE_MODE		},
	{ termbox.KeyArrowRight	,	0	,	ACTION_TURN_CONVERSATION_VIEW_MODE	},
	{ 0						,  ':'	,	ACTION_TURN_COMMAND_MODE			}, /* TODO: need ModAlt field */
	{ termbox.KeyCtrlQ		,	0	,	ACTION_QUIT							},

	{ termbox.KeyArrowUp	,	0	,	ACTION_NEXT_TWEET					},
	{ termbox.KeyArrowDown	,	0	,	ACTION_PREVIOUS_TWEET				},
	{ termbox.KeyHome		,	0	,	ACTION_MOVE_TO_TOP_TWEET			},
	{ termbox.KeyPgup		,	0	,	ACTION_MOVE_TO_TOP_TWEET			},
	{ termbox.KeyEnd		,	0	,	ACTION_MOVE_TO_BOTTOM_TWEET			},
	{ termbox.KeyPgdn		,	0	,	ACTION_MOVE_TO_BOTTOM_TWEET			},
}

/* home timeline action list */
const (
	ACTION_LOAD_PREVIOUSE_TWEETS = iota
	ACTION_LOAD_NEW_TWEETS
)
var homeTimelineKeybindList = []keybind {
	{ termbox.KeySpace,	0,	ACTION_LOAD_PREVIOUSE_TWEETS	},
	{ termbox.KeyCtrlR,	0,	ACTION_LOAD_NEW_TWEETS			},
}

/* input mode action list */
const (
	ACTION_MOVE_LEFT= iota
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
)
var inputModeKeybindList = []keybind {
	{ termbox.KeyArrowLeft	,	0,	ACTION_MOVE_LEFT			},
	{ termbox.KeyArrowRight	,	0,	ACTION_MOVE_RIGHT			},
	{ termbox.KeyArrowUp	,	0,	ACTION_MOVE_UP				},
	{ termbox.KeyArrowDown	,	0,	ACTION_MOVE_DOWN			},
	{ termbox.KeySpace		,	0,	ACTION_INSERT_SPACE			},
	{ termbox.KeyEsc		,	0,	ACTION_EXIT_INPUT_MODE		},
	{ termbox.KeyCtrlG		,	0,	ACTION_EXIT_INPUT_MODE		},
	{ termbox.KeyBackspace	,	0,	ACTION_DELETE_RUNE			},
	{ termbox.KeyBackspace2	,	0,	ACTION_DELETE_RUNE			},
	{ termbox.KeyCtrlA		,	0,	ACTION_MOVE_LINE_TOP		},
	{ termbox.KeyCtrlE		,	0,	ACTION_MOVE_LINE_BOTTOM		},
	{ termbox.KeyCtrlJ		,	0,	ACTION_TURN_CONFIRM_MODE	},
	{ termbox.KeyEnter		,	0,	ACTION_INSERT_NEW_LINE		},
}

/* confirm mode action list */
const (
	ACTION_CANCEL_SUBMIT = iota
	ACTION_SUBMIT_TWEET
)
var confirmModeKeybindList = []keybind {
	{ termbox.KeyEsc,	0	,	ACTION_CANCEL_SUBMIT	},
	{ termbox.KeyCtrlG,	0	,	ACTION_CANCEL_SUBMIT	},
	{ termbox.KeyEnter,	0	,	ACTION_SUBMIT_TWEET		},
}

/* mention view mode action list */
const (
	ACTION_LOAD_PREVIOUSE_MENTIONS = iota
	ACTION_LOAD_NEW_MENTIONS
)
var mentionViewModeKeybindList = []keybind {
	{ termbox.KeyEnter,	0	,	ACTION_LOAD_PREVIOUSE_MENTIONS	},
	{ termbox.KeySpace,	0	,	ACTION_LOAD_PREVIOUSE_MENTIONS	},
	{ termbox.KeyCtrlR,	0	,	ACTION_LOAD_NEW_MENTIONS		},
}

/* conversation mode action list */
const (
	ACTION_EXIT_CONVERSATION_MODE = iota
)
var conversationModeKeybindList = []keybind {
	{ termbox.KeyArrowLeft,		0	,	ACTION_EXIT_CONVERSATION_MODE	},
}


/* user timeline mode action list */
const (
	ACTION_LOAD_PREVIOUSE_USER_TWEETS = iota
	ACTION_LOAD_NEW_USER_TWEETS
)
var userTimelineModeKeybindList = []keybind {
	{ termbox.KeyEnter,	0,	ACTION_LOAD_PREVIOUSE_USER_TWEETS	},
	{ termbox.KeySpace,	0,	ACTION_LOAD_PREVIOUSE_USER_TWEETS	},
	{ termbox.KeyCtrlR,	0,	ACTION_LOAD_NEW_USER_TWEETS			},
}

/* list mode actin list */
const (
	ACTION_LOAD_PREVIOUSE_LIST = iota
	ACTION_LOAD_NEW_LIST
)
var listModeKeybindList = []keybind {
	{ termbox.KeyEnter,	0,	ACTION_LOAD_PREVIOUSE_LIST	},
	{ termbox.KeySpace,	0,	ACTION_LOAD_PREVIOUSE_LIST	},
	{ termbox.KeyCtrlR,	0,	ACTION_LOAD_NEW_LIST		},
}

func (view *view) handleAction(ev termbox.Event, mode KeybindMode) (Action) {
	var action Action = NO_ACTION
	var keybindList []keybind
	switch mode {
		case KEYBIND_MODE_COMMON :
			keybindList = commonKeybindList
		case KEYBIND_MODE_CONVERSATION :
			keybindList = conversationModeKeybindList
		case KEYBIND_MODE_HOME_TIMELINE :
			keybindList = homeTimelineKeybindList
		case KEYBIND_MODE_INPUT :
			keybindList = inputModeKeybindList
		case KEYBIND_MODE_CONFIRM :
			keybindList = confirmModeKeybindList
		case KEYBIND_MODE_MENTION_VIEW :
			keybindList = mentionViewModeKeybindList
		case KEYBIND_MODE_USER_TIMELINE :
			keybindList = userTimelineModeKeybindList
		case KEYBIND_MODE_LIST_VIEW :
			keybindList = listModeKeybindList
	}
	for i := 0; i<len(keybindList); i++ { /* TODO: To handle ModAlt */
		if ev.Key == 0 {	/* kind of CTRL			*/
			if keybindList[i].Ch == ev.Ch {
				action = keybindList[i].Action
				break
			}
		} else {			/* kind of Charactor	*/
			if keybindList[i].Key == ev.Key {
				action = keybindList[i].Action
				break
			}
		}
	}
	return action
}

