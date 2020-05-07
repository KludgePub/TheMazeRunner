package maze

import (
	"math/rand"
	"testing"
	"time"

	"github.com/LinMAD/TheMazeRunnerServer/maze/asset"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestNewMaze(t *testing.T) {
	m := NewMaze(1, 1)

	if len(m.Container) != 1*1 {
		t.Error("Maze map must have Container rows * cols")
	}

	for row, horizonWalls := range m.Walls.Horizontal {
		for _, h := range horizonWalls {
			if h != asset.HorizontalWall {
				t.Error("Initial maze must have all Horizontal Walls")
			}
		}

		for _, verticalWalls := range m.Walls.Vertical[row] {
			if verticalWalls != asset.VerticalWall {
				t.Error("Initial maze must have all Vertical Walls")
			}
		}
	}
}

func TestNewMaze_2x2(t *testing.T) {
	m := NewMaze(2, 2)
	m.Generate()

	mazeStr := PrintMaze(m)
	if len(mazeStr) == 0 {
		t.Error("Maze is empty")
	}

	t.Logf("\nMaze 2x2: \n%s", mazeStr)
}

func TestNewMaze_5x5(t *testing.T) {
	m := NewMaze(5, 5)
	m.Generate()

	mazeStr := PrintMaze(m)
	if len(mazeStr) == 0 {
		t.Error("Maze is empty")
	}

	t.Logf("\nMaze 5x5: \n%s", mazeStr)
}
