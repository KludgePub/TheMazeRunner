package structure

import "fmt"

// Value is an generic
type Value interface{}

// Vertex a single vertex of the graph
type Vertex struct {
	Value Value
}

// String composes vertex Value to string
func (n *Vertex) String() string {
	return fmt.Sprintf("%v", n.Value)
}
