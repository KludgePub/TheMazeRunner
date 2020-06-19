package validator

import (
	"github.com/LinMAD/TheMazeRunnerServer/maze"
)

// GetSolvedPath show path in map
func GetSolvedPath(m maze.Map, from, to maze.Point) []maze.Point {
	stackPath := []maze.Point{from}

	g := maze.DispatchToGraph(&m)

	var isCanMove func(g *maze.Graph, cNode *maze.Node, endPoint maze.Point) bool

	isCanMove = func(g *maze.Graph, cNode *maze.Node, endPoint maze.Point) bool {
		if cNode.Visited {
			return false
		} else if cNode.Point == endPoint {
			return true
		}

		cNode.Visited = true

		if cNode.IsTopNeighbor {
			stackPath = append(stackPath, cNode.TopNeighbor.Point)

			if isCanMove(g, cNode.TopNeighbor, to) {
				return true
			}

			stackPath = stackPath[:len(stackPath)-1]
		}

		if cNode.IsBottomNeighbor {
			stackPath = append(stackPath, cNode.BottomNeighbor.Point)

			if isCanMove(g, cNode.BottomNeighbor, to) {
				return true
			}

			stackPath = stackPath[:len(stackPath)-1]
		}

		if cNode.IsRightNeighbor {
			stackPath = append(stackPath, cNode.RightNeighbor.Point)

			if isCanMove(g, cNode.RightNeighbor, to) {
				return true
			}

			stackPath = stackPath[:len(stackPath)-1]
		}

		if cNode.IsLeftNeighbor {
			stackPath = append(stackPath, cNode.LeftNeighbor.Point)

			if isCanMove(g, cNode.LeftNeighbor, to) {
				return true
			}

			stackPath = stackPath[:len(stackPath)-1]
		}

		return false
	}

	for _, n := range g.Nodes {
		if n.Point == from {
			isCanMove(g, n, to)
			break
		}
	}

	return stackPath
}

// GetPossiblePath from given path
func GetPossiblePath(givenPath []maze.Point, g *maze.Graph) []maze.Point {
	possiblePath := make([]maze.Point, 0)

	for i := 0; i < len(givenPath); i++ {
		var fromNode *maze.Node
		fromPoint := givenPath[i]
		toPoint := fromPoint

		// Get node fromPoint graph

		for _, n := range g.Nodes {
			if n.Point.X == fromPoint.X && n.Point.Y == fromPoint.Y {
				fromNode = n
				break
			}
		}

		// If not found, then given point not possible in maze
		if fromNode == nil {
			return possiblePath
		}

		if len(givenPath) > i+1 {
			toPoint = givenPath[i+1]
		}

		// Check if fromPoint point to point move possible
		if !isPossibleToPass(fromNode, toPoint) {
			possiblePath = append(possiblePath, fromPoint)
			return possiblePath
		}

		possiblePath = append(possiblePath, fromPoint)
	}

	return possiblePath
}

// IsPathPossible validate if given path by points is possible in maze
func IsPathPossible(path []maze.Point, g *maze.Graph) bool {
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
		for i := 0; i < len(path); i++ {
			var to maze.Point
			if i+1 >= len(path) {
				to = path[i]
			} else {
				to = path[i+1]
			}

			if isPossibleToPass(g.Nodes[path[i].GetId()], to) {
				continue
			}

			return false
		}
	}

	return true
}

// canMove validate movement
func isPossibleToPass(from *maze.Node, to maze.Point) bool {
	if from.Point.X == to.X && from.Point.Y == to.Y {
		return true
	}

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
