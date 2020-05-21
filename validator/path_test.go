package validator

import (
	"math/rand"
	"testing"
	"time"

	"github.com/LinMAD/TheMazeRunnerServer/maze"
	"github.com/LinMAD/TheMazeRunnerServer/maze/asset"
)

var (
	m *maze.Map
	g *maze.Graph
)

func prepare(t *testing.T) {
	m = maze.NewMaze(2, 2)

	/* Test maze
	+---+---+
	| S     |
	+   +   +
	| K | E |
	+---+---+
	*/
	m.Entrance = maze.Point{X: 0, Y: 0}
	m.Key = maze.Point{X: 0, Y: 1}
	m.Exit = maze.Point{X: 1, Y: 1}

	m.Container[0][0] = asset.StartingPoint
	m.Container[0][1] = asset.EmptySpace
	m.Container[1][0] = asset.KeyPoint
	m.Container[1][1] = asset.EndingPoint

	m.Walls.Vertical[0][1] = asset.EmptySpace
	m.Walls.Horizontal[1][0] = asset.EmptySpace
	m.Walls.Horizontal[1][1] = asset.EmptySpace

	g = maze.DispatchToGraph(m)


	t.Logf("\n%s", maze.PrintMaze(m))
}

func TestIsPathPossible_EmptyOrOneMove(t *testing.T) {
	prepare(t)

	// Empty path
	if IsPathPossible(make([]maze.Point, 0), g) == true {
		t.Error("It's impossible to findPath if moving path is empty")
	}
	// With one step
	if IsPathPossible([]maze.Point{{0, 1}}, g) == true {
		t.Error("One step move is impossible, must be at least 2 (from and to)")
	}
}
func TestIsPathPossible_WallBangMove(t *testing.T) {
	prepare(t)

	path := make([]maze.Point, 2)

	for _, n := range g.Nodes {
		if n.Entity == asset.KeyPoint {
			path[0] = n.Point
		}
		if n.Entity == asset.EndingPoint {
			path[1] = n.Point
		}
	}

	// Check incorrect path from start to exit (must hit wall)
	if IsPathPossible(path, g) == true {
		t.Error("Move must be impossible due wall hit: K: => E")
	}
}

func TestIsPathPossible_MoveTwoSteps(t *testing.T) {
	prepare(t)

	// Get starting point of maze
	path := make([]maze.Point, 2)
	for _, n := range g.Nodes {
		if n.Entity == asset.StartingPoint {
			path[0] = n.Point
		}
		if n.Entity == asset.KeyPoint {
			path[1] = n.Point
		}
	}

	// Check path to key
	if IsPathPossible(path, g) == false {
		t.Error("Move must be possible: S => K")
	}
}

func TestIsPathPossible_MoveWithTurns(t *testing.T) {
	prepare(t)

	path := make([]maze.Point, 4)

	for _, n := range g.Nodes {
		if n.Entity == asset.KeyPoint {
			path[0] = n.Point
		}
		if n.Entity == asset.StartingPoint {
			path[1] = n.Point
		}
		if n.Entity == asset.EmptySpace {
			path[2] = n.Point
		}
		if n.Entity == asset.EndingPoint {
			path[3] = n.Point
		}
	}

	// Check incorrect path from start to exit (must hit wall)
	if IsPathPossible(path, g) == false {
		t.Error("Move must be possible, corner move: K => S => ' ' => E")
	}
}

func TestSolveMaze_CheckPathToKey(t *testing.T) {
	prepare(t)

	solved := SolvePath(*m, m.Entrance, m.Key, true)
	t.Logf("Path:\n%v", solved)

	if IsPathPossible(solved, g) == false {
		t.Error("Path to key must be found")
	}

	t.Logf("\n%s", maze.PrintMaze(m))
}

func TestSolveMaze_CheckPathToExit(t *testing.T) {
	prepare(t)

	solved := SolvePath(*m, m.Key, m.Exit, true)
	t.Logf("Path:\n%v", solved)

	if IsPathPossible(solved, g) == false {
		t.Error("Path to exit must be found")
	}

	t.Logf("\n%s", maze.PrintMaze(m))
}

func TestSolveMaze_RandomMaze(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	bigMaze := maze.NewMaze(4,4)
	bigMaze.Generate()

	t.Logf("\n%s", maze.PrintMaze(bigMaze))

	t.Logf("S x,y: %v and K x,y: %v", bigMaze.Entrance, bigMaze.Key)
	solved := SolvePath(*bigMaze, bigMaze.Entrance, bigMaze.Key, true)
	t.Logf("Path:%v", solved)

	if IsPathPossible(solved, maze.DispatchToGraph(bigMaze)) == false {
		t.Error("Path to key must be found")
	}

	t.Logf("\n%s", maze.PrintMaze(bigMaze))
}
