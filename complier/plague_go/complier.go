package main

import (
	"fmt"
	"os"
)

func read() string {
	data ,err := os.ReadFile("sample.plague")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func main() {
	//input := "pinkerton \"foo for bar is foo\" (10 + (10 - 6)) for"
	input := read()

	tokens := initTokenizer(input).Tokenize()
	ast := initParser(tokens).parse()
	tast := transformer(ast)
	code := codeGen(tast)
	fmt.Println(code)
}
