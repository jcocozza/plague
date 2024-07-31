package main

import (
	"fmt"
	"strings"
)

type node struct {
	kind       string
	value      string
	name       string
	callee     *node
	expression *node
	body       []node
	params     []node
	arguments  *[]node
	context    *[]node
}

func (n node) StringP() string {
	return n.prettyPrint(0)
}

func (n *node) prettyPrint(indent int) string {
	s := ""
	s += strings.Repeat("  ", indent)

	switch n.kind {
	case "CallExpression":
		s += "(" + n.name + ")\n"
		s += n.params[0].prettyPrint(indent + 1)
		s += n.params[1].prettyPrint(indent + 1)
	default:
		s += n.kind + "\n"
	}
	return s
}

type parser struct {
	loc       int
	pt        []token
	currToken token
}

func initParser(tokens []token) *parser {
	return &parser{
		loc:       0,
		pt:        tokens,
		currToken: tokens[0],
	}
}

func (p *parser) consume() {
	p.loc += 1
	if p.loc < len(p.pt) {
		p.currToken = p.pt[p.loc]
	}
}

func (p *parser) peek() token {
	return p.pt[p.loc+1]
}

func (p *parser) parse() node {
	ast := node{
		kind: "Program",
		body: []node{},
	}
	ast.body = append(ast.body, p.walk())
	return ast
}

func parseToken(tkn token) node {
	if tkn.kind == INT {
		return node{
			kind:  "NumberLiteral",
			value: tkn.value,
		}
	}
	if tkn.IsOperator() {
		return node{
			kind:   "CallExpression",
			name:   tkn.value,
			params: []node{},
		}
	}
	fmt.Println("unhandled token: ", tkn)
	return node{}
}

func (p *parser) walk() node {
	switch {
	case p.currToken.kind == LPAREN || p.currToken.kind == RPAREN: // skip over delimeters - AST's don't need them
		p.consume()
		return p.walk()
	case p.loc < len(p.pt)-1 && p.peek().IsBinaryOperator():
		startNode := parseToken(p.currToken)
		p.consume()
		opNode := parseToken(p.currToken)
		p.consume()
		opNode.params = append(opNode.params, startNode)
		opNode.params = append(opNode.params, p.walk())
		return opNode
	default:
		nd := parseToken(p.currToken)
		p.consume()
		return nd
	}
}

type visitor map[string]func(n *node, p node)

func traverseArray(a []node, p node, v visitor) {
	for _, child := range a {
		traverseNode(child, p, v)
	}
}

func traverser(a node, v visitor) {
	traverseNode(a, node{}, v)
}

func traverseNode(n, p node, v visitor) {
	for k, visitfunc := range v {
		if k == n.kind {
			visitfunc(&n, p)
		}
	}

	switch n.kind {
		case "Program":
			traverseArray(n.body, n, v)
		case "CallExpression":
			traverseArray(n.params, n, v)
		case "NumberLiteral":
			break
		default:
			fmt.Println("unrecognized node type: ", n.kind)
	}
}

var traverseFunc = map[string]func(n *node, p node){
	"NumberLiteral": func(n *node, p node) {
		*p.context = append(*p.context, node{
			kind: "NumberLiteral",
			value: n.value,
		})
	},
	"CallExpression": func(n *node, p node) {
		nn := node{
			kind: "CallExpression",
			callee: &node{
				kind: "Identifier",
				name: n.name,
			},
			arguments: new([]node),
		}
		n.context = nn.arguments

		if p.kind != "CallExpression" {
			es := node{
				kind: "ExpressionStatement",
				expression: &nn,
			}

			*p.context = append(*p.context, es)
		} else {
			*p.context = append(*p.context, nn)
		}
	},
}

func transformer(ast node) node {
	newAst := node{
		kind: "Program",
		body: []node{},
	}
	ast.context = &newAst.body

	traverser(ast, traverseFunc)
	return newAst
}

func codeGen(n node) string {
	switch n.kind {
		case "Program":
			var r []string
			for _, nod := range n.body {
					r = append(r, codeGen(nod))
			}
			return strings.Join(r, "\n")
		case "ExpressionStatement":
			return codeGen(*n.expression)
		case "CallExpression":
			var r []string
			c := codeGen(*n.callee)
			for _, nod := range *n.arguments {
				r = append(r, codeGen(nod))
			}
			rf := strings.Join(r, ", ")
			return c + "(" + rf + ")"
		case "Identifier":
			return n.name
		case "NumberLiteral":
			return n.value
		default:
			fmt.Println("unhandled node type")
			return ""
	}
}
