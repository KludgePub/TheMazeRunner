package maze

import (
	"fmt"
	"testing"
)

func TestDispatch_BlockedMaze(t *testing.T) {
	m := NewMaze(2, 2)

	m.Container[0][0] = keyPoint
	m.Container[0][1] = startingPoint
	m.Container[1][0] = emptySpace
	m.Container[1][1] = endingPoint

	g := Dispatch(m)
	for _, n := range g {
		if n.BottomNeighbor != nil || n.RightNeighbor != nil {
			t.Error("Unexpected neighbor in sealed maze")
			t.Logf("\nMaze %dx%d: \n%s\n", m.Width, m.Height, m)
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
	m.Container[0][0] = keyPoint
	m.Container[0][1] = startingPoint
	m.Container[1][0] = emptySpace
	m.Container[1][1] = endingPoint

	m.Walls.Vertical[0][1] = emptySpace
	m.Walls.Vertical[1][1] = emptySpace
	m.Walls.Horizontal[1][0] = emptySpace

	g := Dispatch(m)

	expected := make([]string, m.Size)
	// current [x, y], right [x, y], bottom [x, y]
	expected[0] = "75 => 0 0 1 0 0 1"     // key
	expected[1] = "83 => 1 0 -1 -1 -1 -1" // start
	expected[2] = "32 => 0 1 1 1 -1 -1"   // empty
	expected[3] = "69 => 1 1 -1 -1 -1 -1" // end

	for _, n := range g {
		isFound := false
		for _, e := range expected {
			str := fmt.Sprint(n)
			if str == e {
				isFound = true
				break
			}
		}

		if !isFound {
			t.Error(fmt.Sprintf("\nUnexpected output from node: %s", fmt.Sprint(n)))
		}
	}

	t.Logf("\nMaze %dx%d: \n%s", m.Width, m.Height, m)
	t.Log("Expected output:")
	for _, e := range expected {
		t.Logf(" %s", e)
	}
	t.Log("Nodes for current, right and bottom locations must match!")
}
