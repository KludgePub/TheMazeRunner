package maze

import "github.com/LinMAD/TheMazeRunnerServer/maze/asset"

// Graph container for graph
type Graph map[Point]*Node

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
	RightNeighbor, BottomNeighbor *Node
}

// Dispatch maze map to related graph
func Dispatch(m *Map) Graph {
	graph := make(Graph)

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			var cNode *Node
			cPoint := Point{X: x, Y: y}

			if n, exist := graph[cPoint]; exist {
				cNode = n
			} else {
				cNode = &Node{
					Entity: m.Container[y][x],
					Point:  cPoint,
				}
			}

			// Check right neighbor
			if m.Width > x+1 && m.Walls.Vertical[y][x+1] != asset.VerticalWall {
				var rNode *Node
				rnp := Point{X: x + 1, Y: y}

				if n, exist := graph[rnp]; exist {
					rNode = n
				} else {
					rNode = &Node{
						Entity: m.Container[y][x+1],
						Point:  rnp,
					}
				}

				cNode.RightNeighbor = rNode
			}

			// Check bottom neighbor
			if m.Height > y+1 && m.Walls.Horizontal[y+1][x] != asset.HorizontalWall {
				var bNode *Node
				bnp := Point{X: x, Y: y + 1}

				if n, exist := graph[bnp]; exist {
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
			graph[cNode.Point] = cNode
		}
	}

	return graph
}