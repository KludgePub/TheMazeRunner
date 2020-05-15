package maze

import "github.com/LinMAD/TheMazeRunnerServer/maze/asset"

// Graph with nodes and paths
type Graph struct {
	// Nodes with their relation
	Nodes map[Point]*Node
}

// Point in maze matrix
type Point struct {
	// X, Y location
	X, Y int
}

// Node represent single cell
type Node struct {
	Entity byte
	// Point holds location
	Point Point
	// RightNeighbor, BottomNeighbor are edged nodes
	LeftNeighbor, RightNeighbor, TopNeighbor, BottomNeighbor *Node
}

// DispatchToGraph maze map to related graph
func DispatchToGraph(m *Map) Graph {
	graph := Graph{
		Nodes:  make(map[Point]*Node, m.Size),
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
			leftX := x-1
			if leftX >= 0 && (leftX == 0 || m.Walls.Vertical[y][leftX] != asset.VerticalWall) {
				var lNode *Node

				lnp := Point{X: leftX, Y: y}

				if n, exist := graph.Nodes[lnp]; exist {
					lNode = n
				} else {
					lNode = &Node{
						Entity: m.Container[y][leftX],
						Point:  lnp,
					}
				}

				cNode.LeftNeighbor = lNode
			}

			// Check right neighbor
			rightX := x+1
			if m.Width > rightX && m.Walls.Vertical[y][rightX] != asset.VerticalWall {
				var rNode *Node

				rnp := Point{X: rightX, Y: y}

				if n, exist := graph.Nodes[rnp]; exist {
					rNode = n
				} else {
					rNode = &Node{
						Entity: m.Container[y][rightX],
						Point:  rnp,
					}
				}

				cNode.RightNeighbor = rNode
			}

			// Check top neighbor
			topY := y-1
			if topY >= 0 && (topY == 0 || m.Walls.Horizontal[topY][x] != asset.HorizontalWall) {
				var topNode *Node

				tnp := Point{X: x, Y: topY}

				if n, exist := graph.Nodes[tnp]; exist {
					topNode = n
				} else {
					topNode = &Node{
						Entity: m.Container[topY][x],
						Point:  tnp,
					}
				}

				cNode.TopNeighbor = topNode
			}

			// Check bottom neighbor
			bottomY := y+1
			if m.Height > bottomY && m.Walls.Horizontal[bottomY][x] != asset.HorizontalWall {
				var bNode *Node

				bnp := Point{X: x, Y: bottomY}

				if n, exist := graph.Nodes[bnp]; exist {
					bNode = n
				} else {
					bNode = &Node{
						Entity: m.Container[bottomY][x],
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
