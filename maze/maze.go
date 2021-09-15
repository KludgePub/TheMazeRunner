package maze

import (
	"bytes"
	"math/rand"

	"github.com/KludgePub/TheMazeRunner/maze/asset"
)

// mazeMap by rows
type Map struct {
	// Container of rows and columns
	Container [][]byte
	// Size of rows and columns
	Size, Height, Width int
	// Walls props
	Walls Walls
	// Key, Entrance, Exit locations
	Key, Entrance, Exit Point
	// KeyCode is unique string for key and exit
	KeyCode string
}

// Walls inside maze map
type Walls struct {
	// Horizontal, Vertical wall, ignore first X
	Horizontal, Vertical [][]byte
}

// Point in maze matrix
type Point struct {
	// X location
	X int `json:"x"`
	// Y location
	Y int `json:"y"`
}

// NewMaze generates a new map
func NewMaze(rows, cols int) *Map {
	// Init maze matrix and fill Walls inside
	m := &Map{
		Size:     rows * cols,
		Height:   rows,
		Width:    cols,
		Key:      Point{},
		Entrance: Point{},
		Exit:     Point{},
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

	m.Entrance.X = rand.Intn(height)
	m.Entrance.Y = rand.Intn(width)
	m.Container[m.Entrance.X][m.Entrance.Y] = asset.StartingPoint

	for {
		eProp := m.Container[m.Exit.X][m.Exit.Y]
		kProp := m.Container[m.Key.X][m.Key.Y]

		if eProp != asset.EndingPoint {
			m.Exit.X = rand.Intn(height)
			m.Exit.Y = rand.Intn(width)

			eProp := m.Container[m.Exit.X][m.Exit.Y]
			if eProp != asset.StartingPoint && eProp != asset.KeyPoint {
				m.Container[m.Exit.X][m.Exit.Y], eProp = asset.EndingPoint, asset.EndingPoint
			}
		}

		if kProp != asset.KeyPoint {
			// TODO Tweak location to be on opposite site from ending point
			m.Key.X = rand.Intn(height)
			m.Key.Y = rand.Intn(width)

			kProp := m.Container[m.Key.X][m.Key.Y]
			if kProp != asset.StartingPoint && kProp != asset.EndingPoint {
				m.Container[m.Key.X][m.Key.Y], kProp = asset.KeyPoint, asset.KeyPoint
			}
		}

		if eProp == asset.EndingPoint && kProp == asset.KeyPoint {
			break
		}
	}
}

// fillMaze will runs recursively to construct maze
func (m *Map) fillMaze(x, y int) {
	m.Container[x][y] = asset.EmptySpace

	for _, direction := range rand.Perm(4) {
		switch direction {
		case asset.Up:
			if x > 0 && m.Container[x-1][y] == 0 {
				m.Walls.Horizontal[x][y] = 0
				m.fillMaze(x-1, y)
			}
		case asset.Down:
			if x < len(m.Container)-1 && m.Container[x+1][y] == 0 {
				m.Walls.Horizontal[x+1][y] = 0
				m.fillMaze(x+1, y)
			}
		case asset.Left:
			if y > 0 && m.Container[x][y-1] == 0 {
				m.Walls.Vertical[x][y] = 0
				m.fillMaze(x, y-1)
			}
		case asset.Right:
			if y < len(m.Container[0])-1 && m.Container[x][y+1] == 0 {
				m.Walls.Vertical[x][y+1] = 0
				m.fillMaze(x, y+1)
			}
		}
	}
}
