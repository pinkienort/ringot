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
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/nsf/termbox-go"
	"os"
	osuser "os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	cli := cli{}
	cli.run()
}

var (
	api          *anaconda.TwitterApi
	stateCh      chan string
	stateClearCh chan int
	tweetmap     *TweetMap
	profilemap   *ProfileMap
	termWidth    int
	termHeight   int
	user         UserConfig
)

type cli struct {
	argShowVersion bool
	argAuthFlag    bool
	argUser        string
}

// UserConfig keeps user configuration, ex) access_token, users name...
type UserConfig struct {
	ID                int64
	ScreenName        string
	UserName          string
	AccessToken       string
	AccessTokenSecret string
}

func (cl *cli) run() {
	cl.parseArgs()
	anaconda.SetConsumerKey(ConsumerKey)
	anaconda.SetConsumerSecret(ConsumerSecret)

	user = cl.setting()

	view := newView()
	stateCh = make(chan string)
	stateClearCh = make(chan int, 2)
	tweetmap = newTweetMap()
	profilemap = newProfileMap()

	if err := termbox.Init(); err != nil {
		fmt.Println("Failed to initialize termbox")
		os.Exit(1)
	}
	defer func() {
		termbox.Close()
		if err := recover(); err != nil {
			panic(err)
		}

	}()
	termbox.SetOutputMode(termbox.Output256)
	termbox.SetInputMode(termbox.InputAlt)

	if os.Getenv("TERM") == "xterm" {
		termbox.SetDisableEscSequence(xtermOffSequences)
	}

	setTermSize(termbox.Size())

	drawText("Now Loading...", 0, 0, ColorWhite, ColorBackground)
	termbox.Flush()

	view.Init()
	view.Loop()

}

func (cl *cli) authorize() (string, string) {
	authURL, tempCre, err := anaconda.AuthorizationURL("")
	if err != nil {
		fmt.Println("Failed to authorize")
		os.Exit(1)
	}
	fmt.Println("Twitter Authorization: " + authURL)
	fmt.Print("PIN:")
	stdinReader := bufio.NewReader(os.Stdin)
	str, err := stdinReader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read stdin")
		os.Exit(1)
	}
	str = strings.TrimRight(str, "\n")
	_, values, err := anaconda.GetCredentials(tempCre, str)
	if err != nil {
		fmt.Println("Failed to authorize")
		os.Exit(1)
	}
	return values.Get("oauth_token"), values.Get("oauth_token_secret")
}

// These const variables are used by setting
const (
	ProfileDir = ".ringot"
	ConfigFile = "config.json"
)

func (cl *cli) setting() UserConfig {
	me, err := osuser.Current()
	if err != nil {
		fmt.Println("Failed to get current user")
		os.Exit(1)
	}
	home := me.HomeDir
	fullpath := filepath.Join(home, ProfileDir, ConfigFile)

	var configSlice []UserConfig
	file, err := os.Open(fullpath)
	notFound := (err != nil)
	if !notFound {
		err = json.NewDecoder(file).Decode(&configSlice)
		if err != nil {
			file.Close()
			fmt.Println("Failed to decode config.json")
			os.Exit(1)
		}
		file.Close()
		if len(configSlice) == 0 || configSlice == nil {
			fmt.Println("Configuration Err")
			os.Exit(1)
		}
	}

	if notFound || cl.argAuthFlag {
		var config UserConfig
		at, ats := cl.authorize()
		config.AccessToken, config.AccessTokenSecret = at, ats
		api = anaconda.NewTwitterApi(at, ats)
		u, err := api.GetSelf(nil)
		if err != nil {
			fmt.Println("Failed to get user profile")
			os.Exit(1)
		}
		config.ID = u.Id
		config.ScreenName = u.ScreenName
		config.UserName = u.Name

		if !notFound {
			exist := false
			for i := range configSlice {
				if configSlice[i].ID == config.ID {
					configSlice[i] = config
					exist = true
					break
				}
			}
			if !exist {
				configSlice = append(configSlice, config)
			}
		} else {
			configSlice = []UserConfig{config}
		}
		binary, err := json.Marshal(configSlice)
		if err != nil {
			fmt.Println("Couldn't marshal configuration")
			os.Exit(1)
		}

		dir := filepath.Join(home, ProfileDir)
		if _, err = os.Stat(dir); err != nil {
			err = os.Mkdir(dir, 0700)
			if err != nil {
				fmt.Println("Couldn't make a directory :" + dir)
				os.Exit(1)
			}
		}
		newFile, err := os.Create(fullpath)
		if err != nil {
			fmt.Println("Couldn't create a config file :" + fullpath)
			os.Exit(1)
		}
		defer newFile.Close()
		var output bytes.Buffer
		json.Indent(&output, binary, "", "\t")
		output.WriteTo(newFile)
		return config
	}
	var config UserConfig
	if len(configSlice) == 1 {
		config = configSlice[0]
	} else {
		if cl.argUser != "" {
			for _, c := range configSlice {
				if cl.argUser == c.ScreenName {
					config = c
				}
			}
			if config.ID == 0 {
				// If argUser is number, accept as index in account list
				num, err := strconv.Atoi(cl.argUser)
				if err == nil && num >= 1 && num <= len(configSlice) {
					config = configSlice[num-1]
				}
			}
		}
		if config.ID == 0 {
			fmt.Println("Account List")
			for i, c := range configSlice {
				s := fmt.Sprintf(" %d) @%s", i+1, c.ScreenName)
				fmt.Println(s)
			}
			stdinReader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("Select:")
				str, err := stdinReader.ReadString('\n')
				if err != nil {
					continue
				}
				str = strings.TrimRight(str, "\n")
				num, err := strconv.Atoi(str)
				if err != nil || num <= 0 || num > len(configSlice) {
					fmt.Println("Invalid!")
				} else {
					config = configSlice[num-1]
					break
				}
			}
		}
	}
	api = anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	api.HttpClient.Timeout = time.Second * 5
	u, err := api.GetSelf(nil)
	if err == nil {
		config.ScreenName = u.ScreenName
		config.UserName = u.Name
	}
	return config
}

func (cl *cli) parseArgs() {
	flag.BoolVar(&cl.argShowVersion, "v", false, "Show version of this software")
	flag.BoolVar(&cl.argAuthFlag, "a", false, "Authorize twitter account")
	flag.StringVar(&cl.argUser, "u", "", "Specify a user")

	flag.Parse()

	if cl.argShowVersion {
		printVersion()
		os.Exit(0)
	}
	if cl.argAuthFlag && cl.argUser != "" {
		fmt.Println("Can't authorize and specify user at once")
		os.Exit(0)
	}
}
