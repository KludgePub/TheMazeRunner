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
	g maze.Graph
)

func init()  {
	m = maze.NewMaze(2, 2)

	/* Test maze
	+---+---+
	| S     |
	+   +   +
	| K | E |
	+---+---+
	*/
	m.Entrance = maze.Point{X: 0, Y: 0}
	m.Key = maze.Point{X: 1, Y: 0}
	m.Exit = maze.Point{X: 1, Y: 1}

	m.Container[0][0] = asset.StartingPoint
	m.Container[0][1] = asset.EmptySpace
	m.Container[1][0] = asset.KeyPoint
	m.Container[1][1] = asset.EndingPoint

	m.Walls.Vertical[0][1] = asset.EmptySpace
	m.Walls.Horizontal[1][0] = asset.EmptySpace
	m.Walls.Horizontal[1][1] = asset.EmptySpace

	g = maze.DispatchToGraph(m)
}

func TestIsPathPossible_EmptyOrOneMove(t *testing.T) {
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
	t.Logf("\n%s", maze.PrintMaze(m))

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
	t.Logf("\n%s", maze.PrintMaze(m))

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
	t.Logf("\n%s", maze.PrintMaze(m))
	solved := SolveMaze(*m)
	if IsPathPossible(solved.ToKey, g) == false {
		t.Errorf("Path to key must be found, given path:\n%v", solved.ToKey)
	}
}

func TestSolveMaze_CheckPathToExit(t *testing.T) {
	t.Logf("\n%s", maze.PrintMaze(m))
	solved := SolveMaze(*m)
	if IsPathPossible(solved.ToExit, g) == false {
		t.Errorf("Path to exit must be found, given path:\n%v", solved.ToExit)
	}
}

func TestSolveMaze_RandomMaze(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	bigMaze := maze.NewMaze(2,2)
	bigMaze.Generate()

	t.Logf("\n%s", maze.PrintMaze(bigMaze))

	solved := SolveMaze(*bigMaze)

	if IsPathPossible(solved.ToKey, g) == false {
		t.Log("Path not found")
		t.Errorf("Path to key:\n %v", solved.ToKey)
	}
	if IsPathPossible(solved.ToExit, g) == false {
		t.Log("Path not found")
		t.Errorf("Path to exit:\n %v", solved.ToExit)
	}

	t.Logf("\n%s", maze.PrintMaze(solved.SolvedMap))
}
