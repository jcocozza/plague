package main

import (
	"fmt"
	"strconv"
)

const (
	// Literal Types
	IntLiteral = iota
	FloatLiteral
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
