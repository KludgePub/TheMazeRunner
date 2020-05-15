package validator

import (
	"testing"

	"github.com/LinMAD/TheMazeRunnerServer/maze"
	"github.com/LinMAD/TheMazeRunnerServer/maze/asset"
)

func TestIsPathPossible(t *testing.T) {
	m := maze.NewMaze(2, 2)
	m.Generate()
	g := maze.DispatchToGraph(m)

	t.Logf("\n%s", maze.PrintMaze(m))

	// Get starting point of maze
	path := make([]maze.Point, 2)
	for _, n := range g.Nodes {
		if n.Entity == asset.StartingPoint {
			path[0] = maze.Point{X: n.Point.X, Y: n.Point.Y}
		}
		if n.Entity == asset.KeyPoint {
			path[1] = maze.Point{X: n.Point.X, Y: n.Point.Y}
		}
	}

	// Check path to key
	if IsPathPossible(path, g) == false {
		t.Error("Move must be possible: S: x0, y1 => K: x0, y0")
		t.Logf("\nMaze map container: %v\nMaze map walls: %v", m.Container, m.Walls)
		t.Logf("\nMovment path: %v", path)
	}

	for _, n := range g.Nodes {
		if n.Entity == asset.StartingPoint {
			path[0] = maze.Point{X: n.Point.X, Y: n.Point.Y}
		}
		if n.Entity == asset.EndingPoint {
			path[1] = maze.Point{X: n.Point.X, Y: n.Point.Y}
		}
	}


	// Check incorrect path from start to exit (must hit wall)
	if IsPathPossible(path, g) == true {
		t.Error("Move must be impossible due wall hit: S: x0, y1 => E: x1, y0")
		t.Logf("\nMaze map container: %v\nMaze map walls: %v", m.Container, m.Walls)
		t.Logf("\nMove path: %v", path)
	}
}