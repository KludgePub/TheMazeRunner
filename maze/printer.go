package maze

import (
	"fmt"

	"github.com/LinMAD/TheMazeRunnerServer/maze/asset"
)

// String implementation to print node
func PrintGraphNode(n *Node) string {
	var x1, y1, x2, y2, x3, y3 = n.Point.X, n.Point.Y, -1, -1, -1, -1

	if n.RightNeighbor != nil {
		x2, y2 = n.RightNeighbor.Point.X, n.RightNeighbor.Point.Y
	}

	if n.BottomNeighbor != nil {
		x3, y3 = n.BottomNeighbor.Point.X, n.BottomNeighbor.Point.Y
	}

	return fmt.Sprintf("%s => %d %d %d %d %d %d", string(n.Entity), x1, y1, x2, y2, x3, y3)
}

// String parsing maze map to text interpretation
func PrintMaze(m *Map) string {
	rightCorner := []byte(fmt.Sprintf("%c\n", asset.Corner))
	rightWall := []byte(fmt.Sprintf("%c\n", asset.VerticalWall))

	var b []byte

	// Make visual map, for each X and column and construct visual relation between sections
	for y, horizonWalls := range m.Walls.Horizontal {

		for _, h := range horizonWalls {
			if h == asset.HorizontalWall || y == 0 {
				b = append(b, asset.HorizontalWallTile...)
			} else {
				b = append(b, asset.HorizontalOpenTile...)
			}
		}

		b = append(b, rightCorner...)

		for x, verticalWalls := range m.Walls.Vertical[y] {
			if verticalWalls == asset.VerticalWall || x == 0 {
				b = append(b, asset.VerticalWallTile...)
			} else {
				b = append(b, asset.VerticalOpenTile...)
			}

			// draw object inside this cell
			if m.Container[y][x] != 0 {
				b[len(b)-2] = m.Container[y][x]
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
