package validator

import (
	"github.com/LinMAD/TheMazeRunnerServer/maze"
)

// SolvedMaze maze and road paths
type SolvedMaze struct {
	// Same map but with path foot prints
	SolvedMap *maze.Map
	// ToKey, ToExit directions from start to key and then to exit
	ToKey, ToExit []maze.Point
}

// SolveMaze show path in map
func SolveMaze(m maze.Map) SolvedMaze {
	solved := SolvedMaze{
		SolvedMap: &m,
		ToKey:     make([]maze.Point, 1),
		ToExit:    make([]maze.Point, 1),
	}

	solved.ToKey = seekForPath(
		solved.SolvedMap,
		maze.Point{X: m.Entrance.X, Y: m.Entrance.Y},
		maze.Point{X: m.Key.X, Y: m.Key.Y},
	)

	solved.ToExit = seekForPath(
		solved.SolvedMap,
		maze.Point{X: m.Key.X, Y: m.Key.Y},
		maze.Point{X: m.Exit.X, Y: m.Exit.Y},
	)

	for _, p := range solved.ToKey {
		solved.SolvedMap.Container[p.X][p.Y] = '.'
	}

	return solved
}

// seekForPath in recursion
func seekForPath(m *maze.Map, from, to maze.Point) []maze.Point {
	stack := []maze.Point{from}

	g := maze.DispatchToGraph(m)

	var isCanMove func(g *maze.Graph, cNode *maze.Node, endPoint maze.Point) bool

	isCanMove = func(g *maze.Graph, cNode *maze.Node, endPoint maze.Point) bool {
		if cNode.Visited {
			return false
		}

		if cNode.Point == endPoint {
			return true
		}

		cNode.Visited = true

		if cNode.TopNeighbor != nil {
			stack = append(stack, cNode.TopNeighbor.Point)

			if !isCanMove(g, g.Nodes[cNode.TopNeighbor.Point], to) {
				stack = stack[:len(stack) - 1]
			}
		}

		if cNode.BottomNeighbor != nil {
			stack = append(stack, cNode.BottomNeighbor.Point)

			if !isCanMove(g, g.Nodes[cNode.BottomNeighbor.Point], to) {
				stack = stack[:len(stack) - 1]
			}
		}

		if cNode.RightNeighbor != nil {
			stack = append(stack, cNode.RightNeighbor.Point)

			if !isCanMove(g, g.Nodes[cNode.RightNeighbor.Point], to) {
				stack = stack[:len(stack) - 1]
			}
		}

		if cNode.LeftNeighbor != nil {
			stack = append(stack, cNode.LeftNeighbor.Point)

			if !isCanMove(g, g.Nodes[cNode.LeftNeighbor.Point], to) {
				stack = stack[:len(stack) - 1]
			}
		}

		return false
	}

	for _, n := range g.Nodes {
		if n.Point == from {
			isCanMove(&g, n, to)
			break
		}
	}

	return stack
}

// IsPathPossible validate if given path by points is possible in maze
func IsPathPossible(path []maze.Point, g maze.Graph) bool {
	if len(path) <= 1 {
		return false
	}

	isPointsValid := true

	for _, p := range path {
		isMovePossible := false
		for _, n := range g.Nodes {
			if n.Point.X == p.X && n.Point.Y == p.Y {
				isMovePossible = true
				break
			}
		}

		if !isMovePossible {
			isPointsValid = false
			break
		}
	}

	if !isPointsValid {
		return false
	}

	if isPointsValid {
		for i := 0; i < len(path); i += 2 {
			var to maze.Point
			from := g.Nodes[path[i]]

			if i+1 >= len(path) {
				to = path[i]
			} else {
				to = path[i+1]
			}

			if isPossibleToPass(from, to) {
				continue
			}

			return false
		}
	}

	return true
}

// canMove validate movement
func isPossibleToPass(from *maze.Node, to maze.Point) (isCanMove bool) {
	// Check if we can move to left
	if from.LeftNeighbor != nil {
		if from.LeftNeighbor.Point.X == to.X && from.LeftNeighbor.Point.Y == to.Y {
			return true
		}
	}

	// Check if we can move to right
	if from.RightNeighbor != nil {
		if from.RightNeighbor.Point.X == to.X && from.RightNeighbor.Point.Y == to.Y {
			return true
		}
	}

	// Check if we can move to top
	if from.TopNeighbor != nil {
		if from.TopNeighbor.Point.X == to.X && from.TopNeighbor.Point.Y == to.Y {
			return true
		}
	}

	// Check if we can move to bottom
	if from.BottomNeighbor != nil {
		if from.BottomNeighbor.Point.X == to.X && from.BottomNeighbor.Point.Y == to.Y {
			return true
		}
	}

	return false
}
