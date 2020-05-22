package maze

import (
	"fmt"

	"github.com/LinMAD/TheMazeRunnerServer/maze/asset"
)

// String implementation to print node
func PrintGraphNode(n *Node, isPrettyPrint bool) string {
	var x1, y1, x2, y2, x3, y3 = n.Point.X, n.Point.Y, -1, -1, -1, -1

	if n.RightNeighbor != nil {
		x2, y2 = n.RightNeighbor.Point.X, n.RightNeighbor.Point.Y
	}

	if n.BottomNeighbor != nil {
		x3, y3 = n.BottomNeighbor.Point.X, n.BottomNeighbor.Point.Y
	}

	if isPrettyPrint {
		return fmt.Sprintf("%s => %d %d %d %d %d %d", string(n.Entity), x1, y1, x2, y2, x3, y3)
	}

	return fmt.Sprintf("%d %d %d %d %d %d", x1, y1, x2, y2, x3, y3)
}

// String parsing maze map to text interpretation
func PrintMaze(m *Map) string {
	rightCorner := []byte(fmt.Sprintf("%c\n", asset.Corner))
	rightWall := []byte(fmt.Sprintf("%c\n", asset.VerticalWall))

	var b []byte

	// Make visual map, for each X and column and construct visual relation between sections
	for x, horizonWalls := range m.Walls.Horizontal {

		for _, h := range horizonWalls {
			if h == asset.HorizontalWall || x == 0 {
				b = append(b, asset.HorizontalWallTile...)
			} else {
				b = append(b, asset.HorizontalOpenTile...)
			}
		}

		b = append(b, rightCorner...)

		for y, verticalWalls := range m.Walls.Vertical[x] {
			if verticalWalls == asset.VerticalWall || y == 0 {
				b = append(b, asset.VerticalWallTile...)
			} else {
				b = append(b, asset.VerticalOpenTile...)
			}

			if m.Container[y][x] == asset.StartingPoint {
				s := Point{X: x, Y: y}
				_ = s
			}

			// draw object inside this cell
			if m.Container[x][y] != 0 {
				b[len(b)-2] = m.Container[x][y]
			}
		}

		b = append(b, rightWall...)
	}

	// End of visual map
	for range m.Walls.Horizontal[0] {
		b = append(b, asset.HorizontalWallTile...)
	}

	b = append(b, rightCorner...)

	return string(b)
}
