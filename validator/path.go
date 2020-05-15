package validator

import "github.com/LinMAD/TheMazeRunnerServer/maze"

// IsPathPossible validate collisions with map
func IsPathPossible(path []maze.Point, g maze.Graph) bool {
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
			from := g.Nodes[maze.Point(path[i])]

			// Check if we can move to left
			if from.LeftNeighbor != nil {
				if from.LeftNeighbor.Point.X == path[i+1].Y && from.LeftNeighbor.Point.Y == path[i+1].Y {
					continue
				}
			}

			// Check if we can move to right
			if from.RightNeighbor != nil {
				if from.RightNeighbor.Point.X == path[i+1].Y && from.RightNeighbor.Point.Y == path[i+1].Y {
					continue
				}
			}

			// Check if we can move to top
			if from.TopNeighbor != nil {
				if from.TopNeighbor.Point.X == path[i+1].Y && from.TopNeighbor.Point.Y == path[i+1].Y {
					continue
				}
			}

			// Check if we can move to bottom
			if from.BottomNeighbor != nil {
				if from.BottomNeighbor.Point.X == path[i+1].Y && from.BottomNeighbor.Point.Y == path[i+1].Y {
					continue
				}
			}

			return false
		}
	}

	return true
}

