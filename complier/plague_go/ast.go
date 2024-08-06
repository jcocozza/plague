package main

import (
	"fmt"
	"strconv"
)

const (
	// Literal Types
	IntLiteral = iota
	FloatLiteral
	StringLiteral
)

type Node interface{}

// node types
type Program struct {
	Body []Node
}

type Identifier struct {
	Name string
}

type Comment string

type Literal struct {
	Type  int
	Value string
}

type NestedExpression struct {
	Body []Node
}

type BinaryOperator struct {
	Op    string
	Left  Node
	Right Node
}

type Arg struct {
	Name string
	Type int
}

type Function struct {
	Name string
	Arguements []Arg
}

type Parser struct {
	loc       int
	tokens    []token
	currToken token
}

func initParser(tokens []token) *Parser {
	return &Parser{
		loc:       0,
		tokens:    tokens,
		currToken: tokens[0],
	}
}

func (p *Parser) consume() {
	p.loc += 1
	if p.loc < len(p.tokens) {
		p.currToken = p.tokens[p.loc]
	}
}

func (p *Parser) peek() token {
	return p.tokens[p.loc+1]
}

func (p *Parser) getNode(binaryStep bool) Node {
	switch {
	case p.currToken.kind == LPAREN:
		nestedNode := NestedExpression{
			Body: []Node{},
		}
		p.consume() // consume '('
		for p.currToken.kind != RPAREN {
			nestedNode.Body = append(nestedNode.Body, p.getNode(false))
		}
		p.consume() // consume the closing ')'
		return nestedNode
	case p.loc < len(p.tokens)-1 && p.peek().IsBinaryOperator() && !binaryStep:
		left := p.getNode(true)
		operatorNode := BinaryOperator{
			Op: p.currToken.value,
		}
		p.consume()
		operatorNode.Left = left
		operatorNode.Right = p.getNode(false)
		return operatorNode
	case p.currToken.kind == FUNC:
		node := Function{
			Name: "",
			Arguements: []Arg{},
		}
		p.consume() // consume 'func'

		funcName := p.currToken.value
		node.Name = funcName
		p.consume() // consume function name

		p.consume() // consume '('
		for p.currToken.kind != RPAREN {
			name := p.currToken.value
			p.consume()
			typ := p.currToken.kind
			p.consume()
			arg := Arg{
				Name: name,
				Type: int(typ),
			}
			fmt.Println(arg)
			node.Arguements = append(node.Arguements, arg)
		}
		p.consume() // consume the ')'
		p.consume() // consume the '{'
		for p.currToken.kind != RBRACE {
			fmt.Println("consuming func content")
			// consume the function here
			p.consume()
		}
		p.consume() // consume the '}'
		return node
	case p.currToken.kind == INT:
		node := Literal{
			Type:  INT,
			Value: p.currToken.value,
		}
		p.consume()
		return node
	case p.currToken.kind == FLOAT:
		node := Literal{
			Type:  FloatLiteral,
			Value: p.currToken.value,
		}
		p.consume()
		return node
	case p.currToken.kind == COMMENT:
		node := Comment(p.currToken.value)
		p.consume()
		return node
	case p.currToken.kind == STRING:
		node := Literal{
			Type:  StringLiteral,
			Value: p.currToken.value,
		}
		p.consume()
		return node
	default:
		fmt.Println("unhandled token")
		fmt.Println("token: ", p.currToken.kind)
		panic("unhandled token " + p.currToken.String() + " " + strconv.Itoa(p.currToken.position))
	}
}

func (p *Parser) walk(binaryStep bool) Node {
	node := p.getNode(binaryStep)
	// we need to check if the next token in a binary operator. If it is then the we need to do left and right nodes
	if p.currToken.IsBinaryOperator() && !binaryStep {
		operatorNode := BinaryOperator{
			Op: p.currToken.value,
		}
		operatorNode.Left = node
		p.consume()
		operatorNode.Right = p.getNode(false) //p.walk(false)
		return operatorNode
	}
	return node
}

func (p *Parser) Parse() Node {
	ast := Program{
		Body: []Node{},
	}
	for p.loc < len(p.tokens) {
		ast.Body = append(ast.Body, p.walk(false))
	}
	return ast
}
