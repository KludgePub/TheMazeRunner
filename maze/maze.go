package maze

import (
	"bytes"
	"math/rand"

	"github.com/LinMAD/TheMazeRunnerServer/maze/asset"
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
	h := bytes.Repeat([]byte{asset.HorizontalWall}, m.Size)
	v := bytes.Repeat([]byte{asset.VerticalWall}, m.Size)

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
	m.Container[m.Entrance.row][m.Entrance.col] = asset.StartingPoint

	for {
		eProp := m.Container[m.Exit.row][m.Exit.col]
		kProp := m.Container[m.Key.row][m.Key.col]

		if eProp != asset.EndingPoint {
			m.Exit.row = rand.Intn(width)
			m.Exit.col = rand.Intn(width)

			eProp := m.Container[m.Exit.row][m.Exit.col]
			if eProp != asset.StartingPoint && eProp != asset.KeyPoint {
				m.Container[m.Exit.row][m.Exit.col], eProp = asset.EndingPoint, asset.EndingPoint
			}
		}

		if kProp != asset.KeyPoint {
			// TODO Tweak location to be on opposite site from ending point
			m.Key.row = rand.Intn(width)
			m.Key.col = rand.Intn(width)

			kProp := m.Container[m.Key.row][m.Key.col]
			if kProp != asset.StartingPoint && kProp != asset.EndingPoint {
				m.Container[m.Key.row][m.Key.col], kProp = asset.KeyPoint, asset.KeyPoint
			}
		}

		if eProp == asset.EndingPoint && kProp == asset.KeyPoint {
			break
		}
	}
}

// fillMaze will runs recursively to construct maze
func (m *Map) fillMaze(startW, startH int) {
	m.Container[startW][startH] = asset.EmptySpace

	for _, direction := range rand.Perm(4) {
		switch direction {
		case asset.Up:
			if startW > 0 && m.Container[startW-1][startH] == 0 {
				m.Walls.Horizontal[startW][startH] = 0
				m.fillMaze(startW-1, startH)
			}
		case asset.Left:
			if startH > 0 && m.Container[startW][startH-1] == 0 {
				m.Walls.Vertical[startW][startH] = 0
				m.fillMaze(startW, startH-1)
			}
		case asset.Down:
			if startW < len(m.Container)-1 && m.Container[startW+1][startH] == 0 {
				m.Walls.Horizontal[startW+1][startH] = 0
				m.fillMaze(startW+1, startH)
			}
		case asset.Right:
			if startH < len(m.Container[0])-1 && m.Container[startW][startH+1] == 0 {
				m.Walls.Vertical[startW][startH+1] = 0
				m.fillMaze(startW, startH+1)
			}
		}
	}
}
