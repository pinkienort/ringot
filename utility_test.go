package main

import (
	"reflect"
	"testing"
)

func TestByteSliceRemove(t *testing.T) {
	src := []byte{0, 1, 2, 4}
	expected := []byte{0}
	val := byteSliceRemove(src, 1, 4)
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v, but %v:", expected, val)
	}

	src = []byte{0}
	expected = []byte{}
	val = byteSliceRemove(src, 0, 1)
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v, but %v:", expected, val)
	}

	src = []byte{0, 1, 2, 3, 4, 5, 6, 7}
	expected = []byte{0, 1, 2, 6, 7}
	val = byteSliceRemove(src, 3, 6)
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v, but %v", expected, val)
	}
}

func TestByteSliceInsert(t *testing.T) {
	src := []byte{0, 1, 2, 4}
	expected := []byte{0, 1, 2, 3, 4}
	val := byteSliceInsert(src, []byte{3}, 3)
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v, but %v", expected, val)
	}

	src = []byte{}
	expected = []byte{0, 1, 2, 3, 4, 5, 6}
	val = byteSliceInsert(src, []byte{0, 1, 2, 3, 4, 5, 6}, 0)
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v, but %v", expected, val)
	}

	src = []byte{0, 1, 2}
	expected = []byte{0, 1, 2, 0, 1, 2, 3, 4, 5, 6}
	val = byteSliceInsert(src, []byte{0, 1, 2, 3, 4, 5, 6}, 3)
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected %v, but %v", expected, val)
	}
}

func TestCenteringStr(t *testing.T) {
	var input, expected, val string

	input = "input"
	expected = "  input  "
	val = centeringStr(input, 9)
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected \"%v\":len(%d), but \"%v\":len(%d)",
			expected, len(expected), val, len(val))
	}

	input = "input"
	expected = "   input  "
	val = centeringStr(input, 10)
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected \"%v\":len(%d), but \"%v\":len(%d)",
			expected, len(expected), val, len(val))
	}
	input = "input"
	expected = "   input   "
	val = centeringStr(input, 11)
	if !reflect.DeepEqual(val, expected) {
		t.Fatalf("Expected \"%v\":len(%d), but \"%v\":len(%d)",
			expected, len(expected), val, len(val))
	}
}

func TestIsScreenNameUsable(t *testing.T) {
	testcase := []rune{'a', 'z', 'A', 'Z', '0', '3', '5', '9', '_'}
	for _, c := range testcase {
		if !isScreenNameUsable(c) {
			t.Fatalf("'%c' is Usable for screen_name", c)
		}
	}

	// Fatal Case
	testcase = []rune{'あ', '/', ':', '[', '{'}
	for _, c := range testcase {
		if isScreenNameUsable(c) {
			t.Fatalf("'%c' is not Usable for screen_name", c)
		}
	}
}

func TestGenerateLabelColorByUserID(t *testing.T) {
	var id int64
	for id = 1; id < 50; id++ {
		lc := generateLabelColorByUserID(id)
		if generateLabelColorByUserID(id) != lc {
			t.Fatalf("generateLabelColorByUserID must return same value when id is same")
		}
		id++
		if id >= 50 {
			break
		}
	}
}
