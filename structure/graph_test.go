package structure

import (
	"testing"
)

func TestGraph_String(t *testing.T) {
	v := &Vertex{Value: "A"}
	v1 := &Vertex{Value: "B"}

	g := Graph{}

	g.AddVertex(v)
	g.AddVertex(v1)

	g.AddEdge(v, v1)

	gs := g.String()

	if "A -> B \nB -> A \n" != gs {
		t.Error("Vertex values or vertex edge is incorrect, expected:\nA -> B \nB -> A \n")
	}
}
