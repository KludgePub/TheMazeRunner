package maze

import (
	"fmt"
	"testing"

	"github.com/LinMAD/TheMazeRunnerServer/maze/asset"
)

func TestDispatch_BlockedMaze(t *testing.T) {
	m := NewMaze(2, 2)

	m.Container[0][0] = asset.KeyPoint
	m.Container[0][1] = asset.StartingPoint
	m.Container[1][0] = asset.EmptySpace
	m.Container[1][1] = asset.EndingPoint

	g := DispatchToGraph(m)
	for _, n := range g.Nodes {
		if n.BottomNeighbor != nil || n.RightNeighbor != nil {
			t.Error("Unexpected neighbor in sealed maze")
			t.Logf("\nMaze %dx%d: \n%s\n", m.Width, m.Height, PrintMaze(m))
			t.Fail()
		}
	}
}

func TestDispatch_SolvableMaze(t *testing.T) {
	m := NewMaze(2, 2)

	/* Test maze
	+---+---+
	| K   S |
	+   +---+
	|     E |
	+---+---+
	*/
	m.Container[0][0] = asset.KeyPoint
	m.Container[0][1] = asset.StartingPoint
	m.Container[1][0] = asset.EmptySpace
	m.Container[1][1] = asset.EndingPoint

	m.Walls.Vertical[0][1] = asset.EmptySpace
	m.Walls.Vertical[1][1] = asset.EmptySpace
	m.Walls.Horizontal[1][0] = asset.EmptySpace

	g := DispatchToGraph(m)

	expected := make([]string, m.Size)
	// current [x, y], right [x, y], bottom [x, y]
	expected[0] = "K => 0 0 1 0 0 1"     // key
	expected[1] = "S => 1 0 -1 -1 -1 -1" // start
	expected[2] = "  => 0 1 1 1 -1 -1"   // empty
	expected[3] = "E => 1 1 -1 -1 -1 -1" // end

	for _, n := range g.Nodes {
		isFound := false
		for _, e := range expected {
			str := PrintGraphNode(n)
			if str == e {
				isFound = true
				break
			}
		}

		if !isFound {
			t.Error(fmt.Sprintf("\nUnexpected output from node: %s", PrintGraphNode(n)))
		}
	}

	t.Logf("\nMaze %dx%d: \n%s", m.Width, m.Height, PrintMaze(m))
	t.Log("Expected output:")
	for _, e := range expected {
		t.Logf(" %s", e)
	}
	t.Log("Nodes for current, right and bottom locations must match!")
}
