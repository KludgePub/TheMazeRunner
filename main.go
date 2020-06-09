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
	row, column := 20, 20

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

	// TCP API handling for game client
	for {
		log.Printf("%s\n", "-> UDP API server, initilizing new connection...")
		conn, cErr := api.NewServerConnection("40", jm)
		if cErr != nil {
			panic(cErr)
		}

		isClosed, handleErr := conn.Handle()
		if handleErr != nil {
			log.Printf("-> UDP API server handling error: %s...", handleErr.Error())
		}
		if isClosed {
			log.Printf("-> TUDP API server is gracefully shutdown...")
			break
		}
	}

	// HTTP API handling for game players
	// TODO add player API
}

// CreateGameWorld maze map
func CreateGameWorld(r, c int) (m *maze.Map, err error) {
	const MaxAttempt = int(^uint(0) >> 1)

	for i := 1; i <= MaxAttempt; i++ {
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
