package maze

import (
	"fmt"
	"testing"

	"github.com/LinMAD/TheMazeRunnerServer/maze/asset"
)

func TestDispatch_HorizontalWall(t *testing.T) {
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
	expected[0] = "K => 0 0 0 1 1 0"     // key
	expected[1] = "S => 0 1 -1 -1 -1 -1" // start
	expected[2] = "  => 1 0 1 1 -1 -1"   // empty
	expected[3] = "E => 1 1 -1 -1 -1 -1" // end

	for _, n := range g.Nodes {
		if n.Entity == asset.EndingPoint {
			if n.TopNeighbor != nil {
				t.Error("Wall must block edge of E => S")
			}
		}
		if n.Entity == asset.StartingPoint {
			if n.BottomNeighbor != nil {
				t.Error("Wall must block edge of S => E")
			}
		}

		isFound := false
		for _, e := range expected {
			str := PrintGraphNode(n, true)
			if str == e {
				isFound = true
				break
			}
		}

		if !isFound {
			t.Error(fmt.Sprintf("\nUnexpected output from node: %s", PrintGraphNode(n, true)))
		}
	}

	t.Logf("\nMaze %dx%d: \n%s", m.Width, m.Height, PrintMaze(m))
	t.Log("Expected output:")
	for _, e := range expected {
		t.Logf(" %s", e)
	}
	t.Log("Nodes for current, right and bottom locations must match!")
}

func TestDispatch_VerticalWall(t *testing.T) {
	m := NewMaze(2, 2)

	/* Test maze
	+---+---+
	| K | S |
	+   +   +
	|     E |
	+---+---+
	*/
	m.Container[0][0] = asset.KeyPoint
	m.Container[0][1] = asset.StartingPoint
	m.Container[1][0] = asset.EmptySpace
	m.Container[1][1] = asset.EndingPoint

	m.Walls.Vertical[1][1] = asset.EmptySpace
	m.Walls.Horizontal[1][1] = asset.EmptySpace
	m.Walls.Horizontal[1][0] = asset.EmptySpace

	g := DispatchToGraph(m)

	expected := make([]string, m.Size)
	// current [x, y], right [x, y], bottom [x, y]
	expected[0] = "K => 0 0 -1 -1 1 0"   // key
	expected[1] = "S => 0 1 -1 -1 1 1"   // start
	expected[2] = "  => 1 0 1 1 -1 -1"   // empty
	expected[3] = "E => 1 1 -1 -1 -1 -1" // end

	for _, n := range g.Nodes {
		if n.Entity == asset.KeyPoint {
			if n.RightNeighbor != nil {
				t.Error("Wall must block edge of K => S")
			}
		}
		if n.Entity == asset.StartingPoint {
			if n.LeftNeighbor != nil {
				t.Error("Wall must block edge of S => K")
			}
		}

		isFound := false
		for _, e := range expected {
			str := PrintGraphNode(n, true)
			if str == e {
				isFound = true
				break
			}
		}

		if !isFound {
			t.Error(fmt.Sprintf("\nUnexpected output from node: %s", PrintGraphNode(n, true)))
		}
	}

	t.Logf("\nMaze %dx%d: \n%s", m.Width, m.Height, PrintMaze(m))
	t.Log("Expected output:")
	for _, e := range expected {
		t.Logf(" %s", e)
	}
	t.Log("Nodes for current, right and bottom locations must match!")
}

func TestDispatch_WallsAround(t *testing.T) {
	m := NewMaze(2, 2)

	/* Test maze
	+---+---+
	| K | S |
	+---+---+
	|   | E |
	+---+---+
	*/
	m.Container[0][0] = asset.KeyPoint
	m.Container[0][1] = asset.StartingPoint
	m.Container[1][0] = asset.EmptySpace
	m.Container[1][1] = asset.EndingPoint

	g := DispatchToGraph(m)

	expected := make([]string, m.Size)
	// current [x, y], right [x, y], bottom [x, y]
	expected[0] = "K => 0 0 -1 -1 -1 -1" // key
	expected[1] = "S => 0 1 -1 -1 -1 -1" // start
	expected[2] = "  => 1 0 -1 -1 -1 -1" // empty
	expected[3] = "E => 1 1 -1 -1 -1 -1" // end

	for _, n := range g.Nodes {
		if n.RightNeighbor != nil || n.LeftNeighbor != nil || n.TopNeighbor != nil || n.BottomNeighbor != nil {
			t.Error("Each node blocked by walls, no edges must be included")
		}
		isFound := false
		for _, e := range expected {
			str := PrintGraphNode(n, true)
			if str == e {
				isFound = true
				break
			}
		}

		if !isFound {
			t.Error(fmt.Sprintf("\nUnexpected output from node: %s", PrintGraphNode(n, true)))
		}
	}

	t.Logf("\nMaze %dx%d: \n%s", m.Width, m.Height, PrintMaze(m))
	t.Log("Expected output:")
	for _, e := range expected {
		t.Logf(" %s", e)
	}
	t.Log("Nodes for current, right and bottom locations must match!")
}
