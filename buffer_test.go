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
	"reflect"
	"testing"
)

func initialize() {
	setTermSize(80, 24)
}

func TestBufferCursor(t *testing.T) {
	var val, expected []byte
	initialize()
	buffer := newBuffer()

	buffer.setContent("a")
	expected = []byte("a")
	val = buffer.content
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v, but %v", expected, val)
	}
	buffer.setContent("a")
	buffer.insertRune('あ')
	buffer.insertRune('0')

	expected = []byte("aあ0")
	val = buffer.content
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v(%s), but %v(%s)", expected, string(expected), val, string(val))
	}
	if buffer.cursorX != len(expected) {
		t.Fatalf("buffer.cursorX is wrong, Expected %v, but %v", len(expected), buffer.cursorX)
	}

	buffer.deleteRuneBackward()
	expected = []byte("aあ")
	val = buffer.content
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v(%s), but %v(%s)", expected, string(expected), val, string(val))
	}
	if buffer.cursorX != len(expected) {
		t.Fatalf("buffer.cursorX is wrong, Expected %v, but %v", len(expected), buffer.cursorX)
	}

	buffer.deleteRuneBackward()
	expected = []byte("a")
	val = buffer.content
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v(%s), but %v(%s)", expected, string(expected), val, string(val))
	}
	if buffer.cursorX != len(expected) {
		t.Fatalf("buffer.cursorX is wrong, Expected %v, but %v", len(expected), buffer.cursorX)
	}

	buffer.cursorMoveBackward()
	buffer.deleteRuneBackward()
	expected = []byte("a")
	val = buffer.content
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v(%s), but %v(%s)", expected, string(expected), val, string(val))
	}
	if buffer.cursorX != 0 {
		t.Fatalf("buffer.cursorX is wrong, Expected %v, but %v", len(expected), buffer.cursorX)
	}

	buffer.setContent("あめんぼあかいな、あいうえお")

	for i := 0; i < 5; i++ {
		buffer.cursorMoveBackward()
	}
	buffer.deleteRuneBackward()

	expected = []byte("あめんぼあかいなあいうえお")
	val = buffer.content
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v(%s), but %v(%s)", expected, string(expected), val, string(val))
	}
	expected = []byte("あめんぼあかいな")
	if buffer.cursorX != len(expected) {
		t.Fatalf("buffer.cursorX is wrong, Expected %v, but %v", len(expected), buffer.cursorX)
	}

	for i := 0; i < 2; i++ {
		buffer.cursorMoveForward()
	}
	for i := 0; i < 2; i++ {
		buffer.deleteRuneBackward()
	}
	buffer.insertRune('愛')

	expected = []byte("あめんぼあかいな愛うえお")
	val = buffer.content
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v(%s), but %v(%s)", expected, string(expected), val, string(val))
	}
	expected = []byte("あめんぼあかいな愛")
	if buffer.cursorX != len(expected) {
		t.Fatalf("buffer.cursorX is wrong, Expected %v, but %v", len(expected), buffer.cursorX)
	}

	buffer.cursorMoveToBottom()
	buffer.insertRune('!')
	buffer.cursorMoveBackward()
	buffer.deleteRuneBackward()
	buffer.insertRune('尾')

	expected = []byte("あめんぼあかいな愛うえ尾!")
	val = buffer.content
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v(%s), but %v(%s)", expected, string(expected), val, string(val))
	}
	expected = []byte("あめんぼあかいな愛うえ尾")
	if buffer.cursorX != len(expected) {
		t.Fatalf("buffer.cursorX is wrong, Expected %v, but %v", len(expected), buffer.cursorX)
	}
}
