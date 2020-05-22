package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/LinMAD/TheMazeRunnerServer/maze"
	"github.com/LinMAD/TheMazeRunnerServer/validator"
)

func init() {
	log.Printf("%s\n", "-> Initilizing the maze server...")
	rand.Seed(time.Now().UnixNano())
}

func main() {
	row, column := 5, 5

	log.Printf("-> Generating new maze (%dx%d)...\n", row, column)

	m := CreateGameWorld(row, column)

	log.Println("-> Maze ready...")
	log.Printf("-> Visual map:\n\n")
	fmt.Println(maze.PrintMaze(m))

	log.Println("-> Interpretation of maze for players...")
	log.Println("-> Nodes: Current, right, bottom nodes with: x,y x,y, x,y")
	log.Println("-> Maze nodes lines:")

	for _, n := range maze.DispatchToGraph(m).Nodes {
		fmt.Printf("%s\n", maze.PrintGraphNode(n, true))
	}

	// TODO Execute API for players to remotely control game
	// TODO Add storage to register: Player, Maze, Score, Locations
}

// CreateGameWorld maze map
func CreateGameWorld(r, c int) (m *maze.Map) {
	i := 0
	for {
		i++
		log.Printf("-> Assemble a maze in (%d) attempt...\n", i)

		m = maze.NewMaze(r, c)
		m.Generate()
		g := maze.DispatchToGraph(m)

		toKey := validator.SolvePath(*m, m.Entrance, m.Key)
		toExit := validator.SolvePath(*m, m.Entrance, m.Exit)

		if validator.IsPathPossible(toKey, g) && validator.IsPathPossible(toExit, g) {
			return
		}
	}
}
