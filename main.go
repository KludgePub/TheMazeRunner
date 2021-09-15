package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KludgePub/TheMazeRunner/api/game"
	"github.com/KludgePub/TheMazeRunner/api/player"
	"github.com/KludgePub/TheMazeRunner/generator"
	"github.com/KludgePub/TheMazeRunner/manager"
	"github.com/KludgePub/TheMazeRunner/maze"
	"github.com/KludgePub/TheMazeRunner/validator"
)

var mazeSize uint

func init() {
	flag.UintVar(&mazeSize, "mazeSize", 20, "Set maze size by squares, default 20x20")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	quit := make(chan int)
	row, column := mazeSize, mazeSize

	log.Printf("-> Generating new maze (%dx%d)...\n", row, column)

	m, mErr := CreateGameWorld(int(row), int(column))
	if mErr != nil {
		log.Fatalf("-> Unable to create game world: %v", mErr)
	}

	log.Println("-> Maze ready...")
	m.KeyCode, mErr = generator.CreateUUID()
	if mErr != nil {
		log.Fatalf("-> Unable to generate key UID: %v", mErr)
	}

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

	go ExecuteServerHTTP(m, mg, gm, h, 8042)
	go ExecuteServerUDP(gm, jm)

	select {
	case <-quit:
		os.Exit(0)
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

		toKey := validator.GetSolvedPath(*m, m.Entrance, m.Key)
		toExit := validator.GetSolvedPath(*m, m.Entrance, m.Exit)

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
func ExecuteServerHTTP(mazeMap *maze.Map, mazeGraph *maze.Graph, gm *manager.GameManager, hostname string, port int) {
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
		}
	}()

	log.Printf("%s %v...\n", "-> HTTP API executor: ready to listen on port", port)
	if err := a.Start(port); err != http.ErrServerClosed {
		panic(fmt.Sprintf("%s %v", "-> HTTP API executor: server failed,", err))
	}
}
