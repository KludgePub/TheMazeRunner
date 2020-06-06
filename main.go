package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/LinMAD/TheMazeRunnerServer/api"
	"github.com/LinMAD/TheMazeRunnerServer/maze"
	"github.com/LinMAD/TheMazeRunnerServer/validator"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	row, column := 5, 5

	log.Printf("-> Generating new maze (%dx%d)...\n", row, column)

	m, mErr := CreateGameWorld(row, column)
	if mErr != nil {
		panic(mErr)
	}

	log.Println("-> Maze ready...")
	log.Printf("-> Visual map:\n")
	log.Printf("\n%s", maze.PrintMaze(m))

	jm, jmErr := json.Marshal(maze.DispatchToGraph(m))
	if jmErr != nil {
		panic(jmErr)
	}

	log.Printf("%s\n", "-> Initilizing the maze server...")
	server, serverErr := api.NewServer("40", jm)
	if serverErr != nil {
		panic(serverErr)
	}

	log.Printf("%s\n", "-> Server ready to handle TCP requests...")
	handleErr := server.Handle()
	if handleErr != nil {
		panic(handleErr)
	}
	// TODO Add storage to register: Player, Maze, Score, Locations
}

// CreateGameWorld maze map
func CreateGameWorld(r, c int) (m *maze.Map, err error) {
	const MaxAttempt = int(^uint(0) >> 1)

	for i := 0; i <= MaxAttempt; i++ {
		log.Printf("-> Assemble a maze in (%d) attempt...\n", i)

		m = maze.NewMaze(r, c)
		m.Generate()
		g := maze.DispatchToGraph(m)

		toKey := validator.SolvePath(*m, m.Entrance, m.Key)
		toExit := validator.SolvePath(*m, m.Entrance, m.Exit)

		if validator.IsPathPossible(toKey, g) && validator.IsPathPossible(toExit, g) {
			return m, nil
		}
	}

	return nil, fmt.Errorf("failed to generate game world, max attemntps reached")
}
