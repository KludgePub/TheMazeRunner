package maze

import (
	"bytes"
	"fmt"
	"math/rand"
)

// MazeMap by rows
type Map struct {
	// Container of rows and columns
	Container [][]byte
	// Size of rows and columns
	Size, Height, Width int
	// Walls props
	Walls Walls
	// Key, Entrance, Exit locations
	Key, Entrance, Exit location
}

// Walls inside maze map
type Walls struct {
	// Horizontal, Vertical wall, ignore first row
	Horizontal, Vertical [][]byte
}

// location shows object position
type location struct {
	// row, col in MazeMap.Container[row][col]
	row, col int
}

// NewMaze generates a new map
func NewMaze(rows, cols int) *Map {
	// Init maze matrix and fill Walls inside
	m := &Map{
		Size:     rows * cols,
		Height:   cols,
		Width:    rows,
		Key:      location{},
		Entrance: location{},
		Exit:     location{},
	}

	c := make([]byte, m.Size)
	h := bytes.Repeat([]byte{horizontalWall}, m.Size)
	v := bytes.Repeat([]byte{verticalWall}, m.Size)

	c2 := make([][]byte, rows)
	h2 := make([][]byte, rows)
	v2 := make([][]byte, rows)

	for i := range h2 {
		c2[i] = c[i*cols : (i+1)*cols]
		h2[i] = h[i*cols : (i+1)*cols]
		v2[i] = v[i*cols : (i+1)*cols]
	}

	m.Container = c2
	m.Walls = Walls{Horizontal: h2, Vertical: v2}

	return m
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
			m.Exit.row = rand.Intn(width)
			m.Exit.col = rand.Intn(width)

			eProp := m.Container[m.Exit.row][m.Exit.col]
			if eProp != startingPoint && eProp != keyPoint {
				m.Container[m.Exit.row][m.Exit.col], eProp = endingPoint, endingPoint
			}
		}

		if kProp != keyPoint {
			// TODO Tweak location to be on opposite site from ending point
			m.Key.row = rand.Intn(width)
			m.Key.col = rand.Intn(width)

			kProp := m.Container[m.Key.row][m.Key.col]
			if kProp != startingPoint && kProp != endingPoint {
				m.Container[m.Key.row][m.Key.col], kProp = keyPoint, keyPoint
			}
		}

		if eProp == endingPoint && kProp == keyPoint {
			break
		}
	}
}

// fillMaze will runs recursively to construct maze
func (m *Map) fillMaze(startW, startH int) {
	m.Container[startW][startH] = emptySpace

	for _, direction := range rand.Perm(4) {
		switch direction {
		case up:
			if startW > 0 && m.Container[startW-1][startH] == 0 {
				m.Walls.Horizontal[startW][startH] = 0
				m.fillMaze(startW-1, startH)
			}
		case left:
			if startH > 0 && m.Container[startW][startH-1] == 0 {
				m.Walls.Vertical[startW][startH] = 0
				m.fillMaze(startW, startH-1)
			}
		case down:
			if startW < len(m.Container)-1 && m.Container[startW+1][startH] == 0 {
				m.Walls.Horizontal[startW+1][startH] = 0
				m.fillMaze(startW+1, startH)
			}
		case right:
			if startH < len(m.Container[0])-1 && m.Container[startW][startH+1] == 0 {
				m.Walls.Vertical[startW][startH+1] = 0
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

	// Make visual map, for each row and column and construct visual relation between sections
	for row, horizonWalls := range m.Walls.Horizontal {

		for _, h := range horizonWalls {
			if h == horizontalWall || row == 0 {
				b = append(b, horizontalWallTile...)
			} else {
				b = append(b, horizontalOpenTile...)
			}
		}

		b = append(b, rightCorner...)

		for column, verticalWalls := range m.Walls.Vertical[row] {
			if verticalWalls == verticalWall || column == 0 {
				b = append(b, verticalWallTile...)
			} else {
				b = append(b, verticalOpenTile...)
			}

			// draw object inside this cell
			if m.Container[row][column] != 0 {
				b[len(b)-2] = m.Container[row][column]
			}
		}

		b = append(b, rightWall...)
	}

	// End of visual map
	for range m.Walls.Horizontal[0] {
		b = append(b, horizontalWallTile...)
	}

	b = append(b, rightCorner...)

	return string(b)
}
