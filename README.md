# Ringot
Twitter Client on terminal  

## Features

* Running on terminal
* Keyboard only operation
* Colorful

## Demo

![Demo](doc/pict/demo01.gif)

## Key bindings

### Timeline

|Key|Command|
|:---|:---|
|<kbd>Ctrl-r</kbd>|Reload|
|<kbd>Ctrl-s</kbd>|Select a buffer (you can tweet from this) |
|<kbd>Ctrl-w</kbd>|Select a buffer with *in_reply_to* |
|<kbd>Ctrl-g</kbd>|Universal cancel button |
|<kbd>Ctrl-f</kbd>|Add a tweet to favorites|
|<kbd>Ctrl-v</kbd>|Retweet a tweet|
|<kbd>Ctrl-o</kbd>|Open a URL with browser|
|<kbd>Ctrl-p</kbd>|Download a picture & Open it|
|<kbd>Home</kbd>|Move cursor to Top|
|<kbd>End</kbd> |Move cursor to Bottom|
|<kbd>PgUp</kbd>|Page Up|
|<kbd>PgDn</kbd>|Page Down|

### Mode
|Key|Command|
|:---|:---|
|<kbd>Ctrl-z</kbd>|Switch to the Home Timeline view |
|<kbd>Ctrl-x</kbd>|Switch to the Mention view |
|<kbd>Ctrl-d</kbd>|Switch to the User Timeline view |
|<kbd>→</kbd>|Show Tweet's conversation |
|<kbd>Alt-x</kbd>|Switch to Command Mode |
|<kbd>Ctrl-q</kbd>|Quit from this application |

### Buffer
|Key|Command|
|:---|:---|
|<kbd>Ctrl-j, Ctrl-Enter</kbd>|Send a tweet |
|<kbd>Ctrl-g</kbd>|Universal cancel button |

### Command
|Command|Operation|
|:---|:---|
|:user *screen_name*|Open a User Timeline |
|:list *list_name*|Open a Twitter List |
|:fav *screen_name*|Open a User favorite Timeline |
|:follow *screen_name*|Follow a user|
|:unfollow *screen_name*|Unfollow a user|
|:set_footer *word*|Set footer for Tweet Edit|
|:unset_footer *word*|Unset footer|

## Installation
Dependencies:  
[go 1.6](https://golang.org/) or newer

```
$ go get github.com/tSU-RooT/ringot
```

### Just want the binary?
Download from releases page.  
Add binary into your $PATH  

## Pull Requests
Bug reports and pull requests are welcome on GitHub

## License
Apache License,Version 2.0  
See LICENSE file.  

```
Ringot


Copyright 2016 tSU-RooT

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

Ringot is using many libraries, list is [doc/library_list.md](doc/library_list.md)  
Library's license information are available in [doc/license_notice.txt](doc/license_notice.txt)  
Some libraries are modified by tSU-RooT  
