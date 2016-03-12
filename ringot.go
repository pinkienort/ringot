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
	"github.com/nsf/termbox-go"
)

// These const variables are used by printVersion
const (
	ApplicationName = "Ringot"
	Version         = 0.1
	CopyrightYear   = "2016"
	Author          = "tSU-RooT"
	NoWarranty      = "Ringot comes with ABSOLUTELY NO WARRANTY"
	LicneseDetail   = "This software is released under the Apache License, Version 2.0"
)

func printVersion() {
	fmt.Printf("%v %v\n", ApplicationName, Version)
	fmt.Printf("Copyright (C) %v %v\n", CopyrightYear, Author)
	fmt.Printf("%v\n", NoWarranty)
	fmt.Printf("%v\n", LicneseDetail)
}

// Configuraion
const (
	CountTweet = 200
)

// Color Configuration
const (
	ColorBackground = termbox.ColorDefault
	ColorRed        = termbox.ColorRed
	ColorWhite      = termbox.ColorWhite
	ColorYellow     = termbox.ColorYellow
	ColorGreen      = termbox.ColorGreen
	ColorBlue       = termbox.ColorBlue
	ColorPink       = termbox.Attribute(214)
	ColorGray1      = termbox.Attribute(0xe9 + 9)
	ColorGray2      = termbox.Attribute(0xe9 + 6)
	ColorGray3      = termbox.Attribute(254)
	ColorLowlight   = termbox.Attribute(240)
)

// Label Colors Configuration
var (
	LabelColors = []termbox.Attribute{
		termbox.ColorRed,
		termbox.ColorGreen,
		termbox.Attribute(40),
		termbox.Attribute(41),
		termbox.Attribute(64),
		termbox.Attribute(66),
		termbox.Attribute(70),
		termbox.Attribute(100),
		termbox.Attribute(110),
		termbox.Attribute(119),
		termbox.Attribute(124),
		termbox.Attribute(126),
		termbox.Attribute(130),
		termbox.Attribute(150),
		termbox.Attribute(160),
		termbox.Attribute(167),
		termbox.Attribute(216),
		termbox.Attribute(227),
	}
	LabelColorMap = make(map[int64]int, 50)
)
