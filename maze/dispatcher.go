package maze

import "github.com/LinMAD/TheMazeRunnerServer/maze/asset"

// Graph with nodes and paths
type Graph struct {
	// Nodes with their relation
	Nodes map[Point]*Node
}

// Node represent single cell
type Node struct {
	// Visited if node was traversed
	Visited bool
	// Entity represent an value
	Entity byte
	// Point holds location
	Point Point
	// RightNeighbor, BottomNeighbor are edged nodes
	LeftNeighbor, RightNeighbor, TopNeighbor, BottomNeighbor *Node
}

// DispatchToGraph assemble graph to provide it to player
func DispatchToGraph(m *Map) *Graph {
	graph := Graph{
		Nodes: make(map[Point]*Node, m.Size),
	}

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			var cNode *Node
			cPoint := Point{X: x, Y: y}

			if n, exist := graph.Nodes[cPoint]; exist {
				cNode = n
			} else {
				cNode = &Node{
					Entity: m.Container[y][x],
					Point:  cPoint,
				}
			}

			// Check left neighbor
			if x-1 >= 0 && (m.Walls.Vertical[y][x] != asset.VerticalWall) {
				var lNode *Node

				lnp := Point{X: x - 1, Y: y}

				if n, exist := graph.Nodes[lnp]; exist {
					lNode = n
				} else {
					lNode = &Node{
						Entity: m.Container[y][x-1],
						Point:  lnp,
					}
				}

				cNode.LeftNeighbor = lNode
			}

			// Check right neighbor
			if m.Width > x+1 && m.Walls.Vertical[y][x+1] != asset.VerticalWall {
				var rNode *Node

				rnp := Point{X: x + 1, Y: y}

				if n, exist := graph.Nodes[rnp]; exist {
					rNode = n
				} else {
					rNode = &Node{
						Entity: m.Container[y][x+1],
						Point:  rnp,
					}
				}

				cNode.RightNeighbor = rNode
			}

			// Check top neighbor
			if y-1 >= 0 && (m.Walls.Horizontal[y][x] != asset.HorizontalWall) {
				var topNode *Node

				tnp := Point{X: x, Y: y - 1}

				if n, exist := graph.Nodes[tnp]; exist {
					topNode = n
				} else {
					topNode = &Node{
						Entity: m.Container[y-1][x],
						Point:  tnp,
					}
				}

				cNode.TopNeighbor = topNode
			}

			// Check bottom neighbor
			if m.Height > y+1 && m.Walls.Horizontal[y+1][x] != asset.HorizontalWall {
				var bNode *Node

				bnp := Point{X: x, Y: y + 1}

				if n, exist := graph.Nodes[bnp]; exist {
					bNode = n
				} else {
					bNode = &Node{
						Entity: m.Container[y+1][x],
						Point:  bnp,
					}
				}

				cNode.BottomNeighbor = bNode
			}

			// Update graph
			graph.Nodes[cNode.Point] = cNode
		}
	}

	return &graph
}
