package main

type Transformer struct{}

func (t *Transformer) Transform(node Node) Node {
	switch n := node.(type) {
	default:
		return n
	}
}
