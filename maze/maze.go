package maze

import (
	"bytes"
	"fmt"
	"math/rand"
)

// Map by rows
type Map struct {
	cells [][]byte
	walls walls
}

type walls struct {
	horizontal [][]byte // ignore first row
	vertical   [][]byte // ignore first of each column
}

// NewMaze generates a new map
func NewMaze(rows, cols int) *Map {
	// Init maze matrix and fill walls inside
	c := make([]byte, rows*cols)
	h := bytes.Repeat([]byte{horizontalWall}, rows*cols)
	v := bytes.Repeat([]byte{verticalWall}, rows*cols)

	c2 := make([][]byte, rows)
	h2 := make([][]byte, rows)
	v2 := make([][]byte, rows)

	for i := range h2 {
		c2[i] = c[i*cols : (i+1)*cols]
		h2[i] = h[i*cols : (i+1)*cols]
		v2[i] = v[i*cols : (i+1)*cols]
	}

	return &Map{c2, walls{horizontal: h2, vertical: v2}}
}

// Generate map
func (m *Map) Generate() {
	m.fillMaze(rand.Intn(len(m.cells)), rand.Intn(len(m.cells[0])))
}

// fillMaze will runs recursively to construct maze
func (m *Map) fillMaze(row, column int) {
	m.cells[row][column] = emptySpace

	for _, direction := range rand.Perm(4) {
		switch direction {
		case up:
			if row > 0 && m.cells[row-1][column] == 0 {
				m.walls.horizontal[row][column] = 0
				m.fillMaze(row-1, column)
			}
		case left:
			if column > 0 && m.cells[row][column-1] == 0 {
				m.walls.vertical[row][column] = 0
				m.fillMaze(row, column-1)
			}
		case down:
			if row < len(m.cells)-1 && m.cells[row+1][column] == 0 {
				m.walls.horizontal[row+1][column] = 0
				m.fillMaze(row+1, column)
			}
		case right:
			if column < len(m.cells[0])-1 && m.cells[row][column+1] == 0 {
				m.walls.vertical[row][column+1] = 0
				m.fillMaze(row, column+1)
			}
		}
	}
}

// String parsing maze map to text interpretation
func (m *Map) String() string {
	rightCorner := []byte(fmt.Sprintf("%c\n", corner))
	rightWall := []byte(fmt.Sprintf("%c\n", verticalWall))

	var b []byte

	for row, horizonWalls := range m.walls.horizontal {
		for _, h := range horizonWalls {
			if h == horizontalWall || row == 0 {
				b = append(b, horizontalWallTile...)
			} else {
				b = append(b, horizontalOpenTile...)
			}
		}

		b = append(b, rightCorner...)

		for column, verticalWalls := range m.walls.vertical[row] {
			if verticalWalls == verticalWall || column == 0 {
				b = append(b, verticalWallTile...)
			} else {
				b = append(b, verticalOpenTile...)
			}

			// draw cell contents
			if m.cells[row][column] != 0 {
				b[len(b)-2] = m.cells[row][column]
			}
		}

		b = append(b, rightWall...)
	}

	for _ = range m.walls.horizontal[0] {
		b = append(b, horizontalWallTile...)
	}

	b = append(b, rightCorner...)

	return string(b)
}
