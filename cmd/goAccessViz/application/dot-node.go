package application

import "gonum.org/v1/gonum/graph"

type dotNode struct {
	graph.Node
	label string
}

func (d *dotNode) DOTID() string {
	return d.label
}
