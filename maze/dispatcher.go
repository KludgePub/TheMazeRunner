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

// TODO REFACTOR: Swap X and Y, by X coordinate we cannot define top and bottom node

// DispatchToGraph maze map to related graph
func DispatchToGraph(m *Map) Graph {
	graph := Graph{
		Nodes: make(map[Point]*Node, m.Size),
	}

	for x := 0; x < m.Height; x++ {
		for y := 0; y < m.Width; y++ {
			var cNode *Node
			cPoint := Point{X: x, Y: y}

			if n, exist := graph.Nodes[cPoint]; exist {
				cNode = n
			} else {
				cNode = &Node{
					Entity: m.Container[x][y],
					Point:  cPoint,
				}
			}

			// Check left neighbor
			leftY := y - 1
			if leftY >= 0 && (leftY == 0 || m.Walls.Vertical[x][leftY] != asset.VerticalWall) {
				var lNode *Node

				lnp := Point{X: x, Y: leftY}

				if n, exist := graph.Nodes[lnp]; exist {
					lNode = n
				} else {
					lNode = &Node{
						Entity: m.Container[x][leftY],
						Point:  lnp,
					}
				}

				cNode.LeftNeighbor = lNode
			}

			// Check right neighbor
			rightY := y + 1
			if m.Width > rightY && m.Walls.Vertical[x][rightY] != asset.VerticalWall {
				var rNode *Node

				rnp := Point{X: x, Y: rightY}

				if n, exist := graph.Nodes[rnp]; exist {
					rNode = n
				} else {
					rNode = &Node{
						Entity: m.Container[x][rightY],
						Point:  rnp,
					}
				}

				cNode.RightNeighbor = rNode
			}

			// Check top neighbor
			topX := x - 1
			if topX >= 0 && (topX == 0 || m.Walls.Horizontal[topX][y] != asset.HorizontalWall) {
				var topNode *Node

				tnp := Point{X: topX, Y: y}

				if n, exist := graph.Nodes[tnp]; exist {
					topNode = n
				} else {
					topNode = &Node{
						Entity: m.Container[topX][y],
						Point:  tnp,
					}
				}

				cNode.TopNeighbor = topNode
			}

			// Check bottom neighbor
			bottomX := x + 1
			if m.Height > bottomX && m.Walls.Horizontal[bottomX][y] != asset.HorizontalWall {
				var bNode *Node

				bnp := Point{X: bottomX, Y: y}

				if n, exist := graph.Nodes[bnp]; exist {
					bNode = n
				} else {
					bNode = &Node{
						Entity: m.Container[bottomX][y],
						Point:  bnp,
					}
				}

				cNode.BottomNeighbor = bNode
			}

			// Update graph
			graph.Nodes[cNode.Point] = cNode
		}
	}

	return graph
}
