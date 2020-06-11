package maze

import (
	"fmt"

	"github.com/LinMAD/TheMazeRunnerServer/maze/asset"
)

// Graph with nodes and paths
type Graph struct {
	// Nodes with their relation
	Nodes map[string]*Node `json:"maze_nodes"`
}

// Node represent single cell
type Node struct {
	// Visited if node was traversed
	Visited bool `json:"-"`
	// Entity represent an value
	Entity byte `json:"entity,omitempty"`
	// Point holds location
	Point Point `json:"point"`

	// IsLeftNeighbor exist
	IsLeftNeighbor bool `json:"is_left_neighbor"`
	// IsRightNeighbor exist
	IsRightNeighbor bool `json:"is_right_neighbor"`
	// IsTopNeighbor exist
	IsTopNeighbor bool `json:"is_top_neighbor"`
	// IsBottomNeighbor exist
	IsBottomNeighbor bool `json:"is_bottom_neighbor"`

	// Do not marshal nodes, it will be recursive

	// LeftNeighbor edged nodes
	LeftNeighbor *Node `json:"-"`
	// RightNeighbor edged nodes
	RightNeighbor *Node `json:"-"`
	// TopNeighbor edged nodes
	TopNeighbor *Node `json:"-"`
	// BottomNeighbor edged nodes
	BottomNeighbor *Node `json:"-"`
}

// GetId creates unique point hash
func (p Point) GetId() string {
	return fmt.Sprintf("x:%d,y:%d", p.X, p.Y)
}

// DispatchToGraph assemble graph to provide it to player
func DispatchToGraph(m *Map) *Graph {
	graph := Graph{
		Nodes: make(map[string]*Node, m.Size),
	}

	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Width; y++ {
			var cNode *Node
			cPoint := Point{X: x, Y: y}

			if n, exist := graph.Nodes[cPoint.GetId()]; exist {
				cNode = n
			} else {
				cNode = &Node{
					Entity: m.Container[x][y],
					Point:  cPoint,
				}
			}

			// Check left neighbor
			if y-1 >= 0 && m.Walls.Vertical[x][y] != asset.VerticalWall {
				var lNode *Node

				lnp := Point{X: x, Y: y - 1}

				if n, exist := graph.Nodes[lnp.GetId()]; exist {
					lNode = n
				} else {
					lNode = &Node{
						Entity: m.Container[x][y-1],
						Point:  lnp,
					}
				}

				cNode.IsLeftNeighbor = true
				graph.Nodes[lNode.Point.GetId()], cNode.LeftNeighbor = lNode, lNode
			}

			// Check right neighbor
			if m.Height > y+1 && m.Walls.Vertical[x][y+1] != asset.VerticalWall {
				var rNode *Node

				rnp := Point{X: x, Y: y + 1}

				if n, exist := graph.Nodes[rnp.GetId()]; exist {
					rNode = n
				} else {
					rNode = &Node{
						Entity: m.Container[x][y+1],
						Point:  rnp,
					}
				}

				cNode.IsRightNeighbor = true
				graph.Nodes[rNode.Point.GetId()], cNode.RightNeighbor = rNode, rNode
			}

			// Check top neighbor
			if x-1 >= 0 && m.Walls.Horizontal[x][y] != asset.HorizontalWall {
				var topNode *Node

				tnp := Point{X: x - 1, Y: y}

				if n, exist := graph.Nodes[tnp.GetId()]; exist {
					topNode = n
				} else {
					topNode = &Node{
						Entity: m.Container[x-1][y],
						Point:  tnp,
					}
				}

				cNode.IsTopNeighbor = true
				graph.Nodes[topNode.Point.GetId()], cNode.TopNeighbor = topNode, topNode
			}

			// Check bottom neighbor
			if m.Width > x+1 && m.Walls.Horizontal[x+1][y] != asset.HorizontalWall {
				var bNode *Node

				bnp := Point{X: x + 1, Y: y}

				if n, exist := graph.Nodes[bnp.GetId()]; exist {
					bNode = n
				} else {
					bNode = &Node{
						Entity: m.Container[x+1][y],
						Point:  bnp,
					}
				}

				cNode.IsBottomNeighbor = true
				graph.Nodes[bNode.Point.GetId()], cNode.BottomNeighbor = bNode, bNode
			}

			// Update graph
			graph.Nodes[cNode.Point.GetId()] = cNode
		}
	}

	return &graph
}
