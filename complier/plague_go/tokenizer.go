package main

import "fmt"

type tokenId int

// tokens
const (
	COMMENT = iota

	LPAREN
	RPAREN
	QUOTE

	operator_begin
	ASSIGN

	binary_operator_begin
	ADD
	SUB
	MULT
	DIV
	binary_operator_end
	operator_end

	// types
	IDENT
	INT
    FLOAT
	STRING

	// keywords
	keyword_start
	FUNC
	keyword_end
)

// a list of tokens indexed by their iota
var tokens = [...]string{
	COMMENT: "COMMENT",
	LPAREN: "(",
	RPAREN: ")",
	QUOTE: "\"",

	ASSIGN: "=",

	ADD: "+",
	SUB: "-",
	MULT: "*",
	DIV: "/",

	IDENT: "IDENT",
	INT: "INT",
	FLOAT: "FLOAT",
	STRING: "STRING",

	FUNC: "func",
}
// keywords is a map of string -> int that maps the string value of a keyword to its iota value
//
// it is initialzed on run and contains all the keywords in between "keyword_start" and "keyword_end"
var keywords map[string]tokenId
func init() {
	keywords = make(map[string]tokenId, keyword_end - (keyword_start + 1))
	for i := keyword_start + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = tokenId(i)
	}
}

// return the keyword iota assignment if the keyword exists
// otherwise, return ident
func keywordLookup(str string) tokenId {
	if tok, exists := keywords[str]; exists {
		return tok
	}
	return IDENT
}

type token struct {
	kind tokenId
	value string
	position int
}

func (t token) IsOperator() bool {
	return (operator_begin < t.kind && t.kind < operator_end)
}

func (t token) IsBinaryOperator() bool {
	return (binary_operator_begin < t.kind && t.kind < binary_operator_end)
}

const (
	LowestPrecedence = 0
)

func (t token) Precedence() int {
	switch t.kind {
	case ADD, SUB:
		return 4
	case MULT, DIV:
		return 5
	}
	return LowestPrecedence
}

type tokenizer struct {
	input string
	runes []rune
	// state
	current int
}

func initTokenizer(input string) *tokenizer {
	return &tokenizer{
		input:   input + "\n", // newline at the end makes thing easier
		runes:   []rune(input),
		current: 0,
	}
}

func (t *tokenizer) consume() {
	t.current++
}

func (t *tokenizer) getChar() string {
	if t.current < len(t.runes) {
		return string(t.runes[t.current])
	}
	return ""
}

func (t *tokenizer) peek() string {
	if t.current < len(t.runes) - 1 {
		return string(t.runes[t.current + 1])
	}
	return ""
}

func (t *tokenizer) Tokenize() []token {
	tokens := []token{}
	for t.current < len(t.runes) {
		ch := t.getChar()
		switch {
		case ch == "=":
			tokens = append(tokens, token{kind: ASSIGN, value: "=", position: t.current})
			t.consume()
		case ch == "(":
			tokens = append(tokens, token{kind: LPAREN, value: "(", position: t.current})
			t.consume()
			continue
		case ch == ")":
			tokens = append(tokens, token{kind: RPAREN, value: ")", position: t.current})
			t.consume()
			continue
		case ch == "+":
			tokens = append(tokens, token{kind: ADD, value: "+", position: t.current})
			t.consume()
			continue
		case ch == "-":
			tokens = append(tokens, token{kind: SUB, value: "-", position: t.current})
			t.consume()
			continue
		case ch == "*":
			tokens = append(tokens, token{kind: MULT, value: "*", position: t.current})
			t.consume()
			continue
		case ch == "/":
			// '/' can be divide or the start of a comment
			switch t.peek() {
			case "/": // single line comment
				comment := ""
				for ch != "\n" {
					comment += ch
					t.consume()
					ch = t.getChar()
				}
				t.consume() // ignore the final \n in the comment
				tokens = append(tokens, token{kind: COMMENT, value: comment, position: t.current})
			case "*":
				t.consume() // consume '/'
				t.consume() // consume '*'
				ch = t.getChar() // start of the contents of the comment
				comment := "/*"
				fmt.Println(comment)
				for { // consume comment until it ends
					if ch == "*" && t.peek() == "/" {
						t.consume() // consume '*'
						t.consume() // conusme '/'
						comment += "*/"
						fmt.Println(comment)
						break
					}
					comment += ch
					t.consume()
					ch = t.getChar()
				}
				tokens = append(tokens, token{kind: COMMENT, value: comment, position: t.current})
			default:
				tokens = append(tokens, token{kind: DIV, value: "/", position: t.current})
				t.consume()
			}
			continue
		case ch == "\"":
			s := "\""
			t.consume()
			ch = t.getChar()
			for ch != "\"" {
				s += ch
				t.consume()
				ch = t.getChar()
			}
			s += "\"" // when we exit the loop, we need to consume the closing '"'
			t.consume()
			tokens = append(tokens, token{kind: STRING, value: s, position: t.current})
			continue
		case isNumber(ch):
			isFloat := false
			numStr := ""
			for isNumber(ch) || ch == "." {
				if ch == "." {
					isFloat = true
				}
				numStr += ch
				t.consume()
				ch = t.getChar()
			}
			if isFloat {
				tokens = append(tokens, token{kind: FLOAT, value: numStr, position: t.current})
			} else {
				tokens = append(tokens, token{kind: INT, value: numStr, position: t.current})
			}
			continue
		case isLetter(ch):
			str := ""
			for isLetter(ch) {
				str += ch
				t.consume()
				ch = t.getChar()
			}
			if len(str) > 1 { // all keywords will have length > 1 so we can avoid lookup in that case
				tok := keywordLookup(str)
				tokens = append(tokens, token{kind: tok, value: str, position: t.current})
			} else {
				tokens = append(tokens, token{kind: IDENT, value: str, position: t.current})
			}
			continue
		default: // ignore everything else
			t.consume()
			continue
		}
	}
	return tokens
}
