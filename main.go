package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LinMAD/TheMazeRunnerServer/api/game"
	"github.com/LinMAD/TheMazeRunnerServer/api/player"
	"github.com/LinMAD/TheMazeRunnerServer/manager"
	"github.com/LinMAD/TheMazeRunnerServer/maze"
	"github.com/LinMAD/TheMazeRunnerServer/validator"
)

var isRunning = true

func init() {
	// TODO fix: if exec all tests this will break static data
	rand.Seed(time.Now().UnixNano())
}

func main() {
	row, column := 2, 2 // TODO Read from input params or json config

	log.Printf("-> Generating new maze (%dx%d)...\n", row, column)

	m, mErr := CreateGameWorld(row, column)
	if mErr != nil {
		panic(mErr)
	}

	log.Println("-> Maze ready...")
	log.Printf("-> Visual map:\n")
	log.Printf("\n%s", maze.PrintMaze(m))

	mg := maze.DispatchToGraph(m)
	jm, jmErr := json.Marshal(mg)
	if jmErr != nil {
		log.Fatalf("-> Unable to marshal game world: %v", jmErr)
	}

	h, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname: %v", err)
	}

	gm := manager.NewGameManager(m)

	go ExecuteServerHTTP(m, mg, gm, h, 80)
	go ExecuteServerUDP(gm, jm)

	for isRunning {
		// TODO When player submitting movement path:
		// TODO 1. validate if it's reachable
		// TODO 2. Move to possible point
		// TODO 3. Save new location
		// TODO 4. Report to game client via UDP about player new location and with correct movement path
	}
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

// ExecuteServerUDP API handling for game client
func ExecuteServerUDP(gm *manager.GameManager, gameMap []byte) {
	for {
		log.Printf("%s\n", "-> UDP API executer: initilizing new connection...")
		conn, cErr := game.NewServerConnection("40", gameMap)
		if cErr != nil {
			panic(cErr)
		}

		isClosed, handleErr := conn.Handle(gm)
		if handleErr != nil {
			log.Printf("-> UDP API executer: handling error: %s...", handleErr.Error())
		}
		if isClosed {
			log.Printf("-> UDP API executer: is gracefully shutdown...")
			break
		}
	}
}

// ExecuteServerHTTP API handling for players
func ExecuteServerHTTP(mazeMap *maze.Map, mazeGraph *maze.Graph, gm *manager.GameManager,hostname string, port int) {
	a := player.NewPlayerApi(gm, mazeMap, mazeGraph, hostname)

	go func() { // shutdown gracefully
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig

		log.Println("-> HTTP API executor: Performing graceful shutdown of HTTP API...")
		select {
		case <-time.After(1 * time.Second):
			if err := a.Shutdown(); err != nil {
				log.Fatalf("-> HTTP API executor: Failed to shutdown server, %v", err)
			}
			isRunning = false
		}
	}()

	log.Printf("%s %v...\n", "-> HTTP API executor: ready to listen on port", port)
	if err := a.Start(port); err != http.ErrServerClosed {
		panic(fmt.Sprintf("%s %v", "-> HTTP API executor: server failed,", err))
	}
}
