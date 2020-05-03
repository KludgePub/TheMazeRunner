package maze

import (
	"bytes"
	"fmt"
	"github.com/LinMAD/TheMazeRunnerServer/maze/search"
	"math/rand"
)

// Map by rows
type Map struct {
	cells    [][]byte
	walls    walls
	Entrance location
	Exit     location
}

type walls struct {
	horizontal [][]byte // ignore first row
	vertical   [][]byte // ignore first of each column
}

// location show object position
type location struct {
	heightPos int
	widthPos  int
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

	return &Map{
		c2,
		walls{horizontal: h2, vertical: v2},
		location{},
		location{},
	}
}

// Generate map
func (m *Map) Generate() {
	height := len(m.cells)
	width := len(m.cells[0])

	m.fillMaze(rand.Intn(height), rand.Intn(width))

	m.Entrance.heightPos = rand.Intn(height)
	m.Entrance.widthPos = rand.Intn(width)
	m.cells[m.Entrance.heightPos][m.Entrance.widthPos] = startingPoint

	for {
		m.Exit.widthPos = rand.Intn(width)
		m.Exit.widthPos = rand.Intn(width)

		if m.cells[m.Exit.heightPos][m.Exit.widthPos] != startingPoint {
			m.cells[m.Exit.heightPos][m.Exit.widthPos] = endingPoint
			break
		}
	}

	search.Solve()
}

// fillMaze will runs recursively to construct maze
func (m *Map) fillMaze(startW, startH int) {
	m.cells[startW][startH] = emptySpace

	for _, direction := range rand.Perm(4) {
		switch direction {
		case up:
			if startW > 0 && m.cells[startW-1][startH] == 0 {
				m.walls.horizontal[startW][startH] = 0
				m.fillMaze(startW-1, startH)
			}
		case left:
			if startH > 0 && m.cells[startW][startH-1] == 0 {
				m.walls.vertical[startW][startH] = 0
				m.fillMaze(startW, startH-1)
			}
		case down:
			if startW < len(m.cells)-1 && m.cells[startW+1][startH] == 0 {
				m.walls.horizontal[startW+1][startH] = 0
				m.fillMaze(startW+1, startH)
			}
		case right:
			if startH < len(m.cells[0])-1 && m.cells[startW][startH+1] == 0 {
				m.walls.vertical[startW][startH+1] = 0
				m.fillMaze(startW, startH+1)
			}
		}
	}
}
//
//func (m *Map) isPossibleMaze() bool {
//	// Check in recursion if path from start to end solvable
//	var move func(row, col, dir int) bool
//
//	// Make a move by height and width in cell map
//	move = func(h, w, direction int) bool {
//		if h == m.Exit.heightPos && w == m.Exit.widthPos {
//			m.cells[h][w] = endingPoint
//			return true
//		}
//
//		if direction != down && m.walls.horizontal[h][w] == 0 {
//			if move(h-1, w, up) {
//				m.cells[h][w] = '^'
//				m.cells[h][w] = '^'
//				return true
//			}
//		}
//
//		if direction != up && h+1 < len(m.walls.horizontal) && m.walls.horizontal[h+1][w] == 0 {
//			if move(h+1, w, down) {
//				m.cells[h][w] = 'v'
//				m.cells[h+1][w] = 'v'
//				return true
//			}
//		}
//
//		if direction != left && w+1 < len(m.walls.vertical[0]) && m.walls.vertical[h][w+1] == 0 {
//			if move(h, w+1, right) {
//				m.cells[h][w] = '>'
//				m.cells[h][w+1] = '>'
//				return true
//			}
//		}
//
//		if direction != right && m.cells[h][w] == 0 {
//			if move(h, w-1, left) {
//				m.cells[h][w] = '<'
//				m.cells[h][w] = '<'
//				return true
//			}
//		}
//
//		return false
//	}
//
//	return move(m.Entrance.heightPos, m.Entrance.widthPos, up)
//}

// String parsing maze map to text interpretation
func (m *Map) String() string {
	fmt.Println(m.cells)
	rightCorner := []byte(fmt.Sprintf("%c\n", corner))
	rightWall := []byte(fmt.Sprintf("%c\n", verticalWall))

	var b []byte

	for row, horizonWalls := range m.walls.horizontal {
		for _, h := range horizonWalls {
			if h == horizontalWall || row == 0 {
				b = append(b, horizontalWallTile...)
			} else {
				b = append(b, horizontalOpenTile...)
				if h != horizontalWall && h != 0 {
					b[len(b)-2] = h
				}
			}
		}

		b = append(b, rightCorner...)

		for column, verticalWalls := range m.walls.vertical[row] {
			if verticalWalls == verticalWall || column == 0 {
				b = append(b, verticalWallTile...)
			} else {
				b = append(b, verticalOpenTile...)
				if verticalWalls != verticalWall && verticalWalls != 0 {
					b[len(b)-4] = verticalWalls
				}
			}

			// draw cell contents
			if m.cells[row][column] != 0 {
				b[len(b)-2] = m.cells[row][column]
			}
		}

		b = append(b, rightWall...)
	}

	for range m.walls.horizontal[0] {
		b = append(b, horizontalWallTile...)
	}

	b = append(b, rightCorner...)

	return string(b)
}
