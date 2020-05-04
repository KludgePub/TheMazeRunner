package maze

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/LinMAD/TheMazeRunnerServer/maze/search"
)

// Map by rows
type Map struct {
	// Container of rows and columns
	Container [][]byte
	// walls props
	walls walls
	// Key need to open exit
	Key location
	// Entrance to the maze
	Entrance location
	// Exit from maze
	Exit location
}

type walls struct {
	horizontal [][]byte // ignore first row
	vertical   [][]byte // ignore first of each column
}

// location shows object position
type location struct {
	// row in Map.Container[row]
	row int
	// col in Map.Container[row][col]
	col int
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
		location{},
	}
}

// Generate map
func (m *Map) Generate() {
	height := len(m.Container)
	width := len(m.Container[0])

	m.fillMaze(rand.Intn(height), rand.Intn(width))

	m.Entrance.row = rand.Intn(height)
	m.Entrance.col = rand.Intn(width)
	m.Container[m.Entrance.row][m.Entrance.col] = startingPoint

	for {
		eProp := m.Container[m.Exit.row][m.Exit.col]
		kProp := m.Container[m.Key.row][m.Key.col]

		if eProp != endingPoint {
			m.Exit.col = rand.Intn(width)
			m.Exit.col = rand.Intn(width)

			if eProp != startingPoint && eProp != keyPoint {
				m.Container[m.Exit.row][m.Exit.col], eProp = endingPoint, endingPoint
			}
		}

		if kProp != endingPoint {
			// TODO Tweak location to be on opposite site from ending point
			m.Key.col = rand.Intn(width)
			m.Key.col = rand.Intn(width)

			if kProp != startingPoint && kProp != endingPoint {
				m.Container[m.Key.row][m.Key.col], kProp = keyPoint, keyPoint
			}
		}

		if eProp == endingPoint && kProp == keyPoint {
			break
		}
	}

	search.Solve()
}

// fillMaze will runs recursively to construct maze
func (m *Map) fillMaze(startW, startH int) {
	m.Container[startW][startH] = emptySpace

	for _, direction := range rand.Perm(4) {
		switch direction {
		case up:
			if startW > 0 && m.Container[startW-1][startH] == 0 {
				m.walls.horizontal[startW][startH] = 0
				m.fillMaze(startW-1, startH)
			}
		case left:
			if startH > 0 && m.Container[startW][startH-1] == 0 {
				m.walls.vertical[startW][startH] = 0
				m.fillMaze(startW, startH-1)
			}
		case down:
			if startW < len(m.Container)-1 && m.Container[startW+1][startH] == 0 {
				m.walls.horizontal[startW+1][startH] = 0
				m.fillMaze(startW+1, startH)
			}
		case right:
			if startH < len(m.Container[0])-1 && m.Container[startW][startH+1] == 0 {
				m.walls.vertical[startW][startH+1] = 0
				m.fillMaze(startW, startH+1)
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
			if m.Container[row][column] != 0 {
				b[len(b)-2] = m.Container[row][column]
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
