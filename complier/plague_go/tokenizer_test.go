package main

import (
	"reflect"
	"testing"
)

func Test_tokenizer_Tokenize(t *testing.T) {
	type fields struct {
		input   string
		runes   []rune
		current int
	}
	tests := []struct {
		name   string
		fields fields
		want   []token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &tokenizer{
				input:   tt.fields.input,
				runes:   tt.fields.runes,
				current: tt.fields.current,
			}
			if got := tr.Tokenize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenizer.Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
