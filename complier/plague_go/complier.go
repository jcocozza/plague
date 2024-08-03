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

func complier(input string) string {
	tokens := initTokenizer(input).Tokenize()
	//fmt.Println("tokens: ", tokens)
	ast := initParser(tokens).parse()
	fmt.Println("ast: ", ast)
	tast := transformer(ast)
	fmt.Println("transformed: ", tast)
	return codeGen(tast)
}

func main() {
	//input := "pinkerton \"foo for bar is foo\" (10 + (10 - 6)) for"
	input := read()
	code := complier(input)
	fmt.Println(code)
}
