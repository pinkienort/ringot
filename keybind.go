package main

import termbox "github.com/nsf/termbox-go"

type Action uint8
type kybd struct {
	Key    termbox.Key
	Ch     rune
	Action Action
}

/* Action list corresponding keybind(kybd) */
const (
	ACT_LIKE_TWEET = iota
	ACT_MENTION
	ACT_RETWEET
	ACT_OPEN_IMAGES
	ACT_NEXT_TWEET
	ACT_PREVIOUS_TWEET
	ACT_PAGE_DOWN
	ACT_PAGE_UP
	ACT_LOAD_NEW_TWEETS
	ACT_GO_INPUT_MODE
	ACT_GO_COMMAND_MODE
	ACT_GO_HOME_TIMELINE_MODE
	ACT_GO_CONVERSATION_VIEW_MODE
	ACT_GO_MENTION_VIEW_MODE
	ACT_GO_USER_TIMELINE_MODE
	ACT_QUIT
	ACT_OPEN_URL
	ACT_SHOW_HELP
	NO_ACTION = 0xff - 1
)

var kybdList = []kybd{
	{	0				,	'l'	,	ACT_LIKE_TWEET					},
	{	0				,	'r'	,	ACT_MENTION						},
	{	0				,	't'	,	ACT_RETWEET						},
	{	0				,	'o'	,	ACT_OPEN_URL					},
	{	0				,	'p'	,	ACT_OPEN_IMAGES					},
	{	0				,	'j'	,	ACT_NEXT_TWEET					},
	{	0				,	'k'	,	ACT_PREVIOUS_TWEET				},
	{	0				,	'd'	,	ACT_PAGE_DOWN					},
	{	0				,	'u'	,	ACT_PAGE_UP						},
	{	0				,	'?'	,	ACT_SHOW_HELP					},
	{	0				,	'.'	,	ACT_LOAD_NEW_TWEETS				},
	{	0				,	'n'	,	ACT_GO_INPUT_MODE				},
	{	0				,	':'	,	ACT_GO_COMMAND_MODE				},
	{	termbox.KeyCtrlZ,	0	,	ACT_GO_HOME_TIMELINE_MODE		},
	{	termbox.KeyCtrlC,	0	,	ACT_GO_CONVERSATION_VIEW_MODE	},
	{	termbox.KeyCtrlR,	0	,	ACT_GO_MENTION_VIEW_MODE		},
	{	termbox.KeyCtrlD,	0	,	ACT_GO_USER_TIMELINE_MODE		},
	{	termbox.KeyCtrlQ,	0	,	ACT_QUIT						},
}

func (view *view) Action(ev termbox.Event) (Action) {
	var action Action = NO_ACTION
	for i := 0; i<len(kybdList); i++ {
		if ev.Key == 0 {	/* kind of CTRL			*/
			if kybdList[i].Ch == ev.Ch {
				action = kybdList[i].Action
				break
			}
		} else {			/* kind of Charactor	*/
			if kybdList[i].Key == ev.Key {
				action = kybdList[i].Action
				break
			}
		}
	}
	return action
}

