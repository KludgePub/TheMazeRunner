package maze

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestNewMaze(t *testing.T) {
	m := NewMaze(1, 1)

	if len(m.Container) != 1*1 {
		t.Error("Maze map must have Container rows * cols")
	}

	for row, horizonWalls := range m.walls.horizontal {
		for _, h := range horizonWalls {
			if h != horizontalWall {
				t.Error("Pure maze must have all horizontal walls")
			}
		}

		for _, verticalWalls := range m.walls.vertical[row] {
			if verticalWalls != verticalWall {
				t.Error("Pure maze must have all vertical walls")
			}
		}
	}
}

func TestNewMaze2x2(t *testing.T) {
	m := NewMaze(2, 2)
	m.Generate()

	mazeStr := m.String()
	if len(mazeStr) == 0 {
		t.Error("Maze is empty")
	}

	t.Logf("\nMaze 2x2: \n%s", mazeStr)
}

func TestNewMaze5x5(t *testing.T) {
	m := NewMaze(5, 5)
	m.Generate()

	mazeStr := m.String()
	if len(mazeStr) == 0 {
		t.Error("Maze is empty")
	}

	t.Logf("\nMaze 5x5: \n%s", mazeStr)
}
