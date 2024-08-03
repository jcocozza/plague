package main

import (
	"fmt"
	"testing"
)

func TestCompiler(t *testing.T) {
	var tests  = []struct {
		input string
		want string
	}{
		//{"",""},
		{"1 + 2", "1 + 2"},
		{"(1 + 2)", "(1 + 2)"},
		{"1 + (1 + 2)", "1 + (1 + 2)"},
		{"(1 + 2) + 3", "(1 + 2) + 3"},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s, %s", tt.input, tt.want)
		t.Run(testname, func(t *testing.T) {
			res := complier(tt.input)
			if res != tt.want {
				t.Errorf("got %s, wanted %s", res, tt.want)
			}
		})
	}
}
