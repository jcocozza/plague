package main

import (
	"fmt"
	"reflect"
)

type CodeGenerator struct{}

func (cg *CodeGenerator) Generate(node Node) string {
	switch n := node.(type) {
	case Program:
		output := ""
		for _, nd := range n.Body {
			output += cg.Generate(nd)
		}
		return output
	case Comment:
		return string(n) + "\n"
	case Literal:
		return n.Value
	case BinaryOperator:
		leftStr := cg.Generate(n.Left)
		rightStr := cg.Generate(n.Right)
		return leftStr + " " + n.Op + " " + rightStr
	case NestedExpression:
		inner := ""
		for _, nd := range n.Body {
			inner += cg.Generate(nd)
		}
		return "(" + inner + ")"
	case Identifier:
		return n.Name
	default:
		fmt.Println("default hit for:", n)
		nodeType := reflect.TypeOf(node)
		fmt.Println("Node type:", nodeType)
		return ""
	}
}
