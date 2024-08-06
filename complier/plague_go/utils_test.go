package main

import (
	"fmt"
	"testing"
)

func TestIsNumber(t *testing.T) {
	var tests = []struct {
		char string
		want bool
	}{
		{"", false},
		{"f", false},
		{"1", true},
		{"2", true},
		{"3", true},
		{"4", true},
		{"5", true},
		{"6", true},
		{"7", true},
		{"8", true},
		{"9", true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s, %t", tt.char, tt.want)
		t.Run(testname, func(t *testing.T) {
			res := isNumber(tt.char)
			if res != tt.want {
				t.Errorf("got %t, wanted %t", res, tt.want)
			}
		})
	}
}

func TestIsLetter(t *testing.T) {
	var tests = []struct {
		char string
		want bool
	}{
		{"", false},
		{"7", false},
		{"8", false},
		{"9", false},
		{"a", true},
		{"b", true},
		{"c", true},
		{"d", true},
		{"e", true},
		{"f", true},
		{"g", true},
		{"h", true},
		{"i", true},
		{"j", true},
		{"A", true},
		{"B", true},
		{"C", true},
		{"D", true},
		{"E", true},
		{"F", true},
		{"G", true},
		{"H", true},
		{"I", true},
		{"J", true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s, %t", tt.char, tt.want)
		t.Run(testname, func(t *testing.T) {
			res := isLetter(tt.char)
			if res != tt.want {
				t.Errorf("got %t, wanted %t", res, tt.want)
			}
		})
	}
}
