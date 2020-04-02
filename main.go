package main

import (
	"fmt"
	"github.com/LinMAD/TheMazeRunnerServer/structure"
)

func main()  {
	fmt.Println("The maze runner")

	v1 := &structure.Vertex{Value: "A"}
	v2 := &structure.Vertex{Value: "B"}
	v3 := &structure.Vertex{Value: "C"}
	v4 := &structure.Vertex{Value: "D"}
	v5 := &structure.Vertex{Value: "E"}

	g := structure.Graph{}

	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)
	g.AddVertex(v4)
	g.AddVertex(v5)

	g.AddEdge(v1, v2)
	g.AddEdge(v2, v4)
	g.AddEdge(v2, v5)
	g.AddEdge(v3, v5)

	fmt.Println(g.String())
}
