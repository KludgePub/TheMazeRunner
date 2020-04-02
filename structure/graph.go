package structure

import (
	"sync"
)

// Graph with vertices and relations
type Graph struct {
	vertices []*Vertex
	edges    map[Vertex][]*Vertex
	lock     sync.RWMutex
}

// AddVertex to graph
func (g *Graph) AddVertex(n *Vertex) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.vertices = append(g.vertices, n)
}

// AddEdge to graph and related vertices
func (g *Graph) AddEdge(n1, n2 *Vertex) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.edges == nil {
		g.edges = make(map[Vertex][]*Vertex)
	}

	g.edges[*n1] = append(g.edges[*n1], n2)
	g.edges[*n2] = append(g.edges[*n2], n1)
}

// String compose graph to string
func (g *Graph) String() (s string) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	for i := 0; i < len(g.vertices); i++ {
		s += g.vertices[i].String() + " -> "
		near := g.edges[*g.vertices[i]]

		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}

		s += "\n"
	}

	return
}
