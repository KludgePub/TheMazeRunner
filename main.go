package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/LinMAD/TheMazeRunnerServer/maze"
)

func init() {
	fmt.Printf("\n%s\n", "-> Initilizing the maze server...")
	rand.Seed(time.Now().UnixNano())
}

func main() {
	fmt.Println("-> Generating new maze...")

	r, c := 2, 2 // Size of maze map
	m := maze.NewMaze(r, c)
	m.Generate()

	fmt.Println("-> Maze ready...")
	fmt.Println("- Bytes map:\n", m.Container)
	fmt.Println("- Visual map: ")
	fmt.Println(m)

	// TODO Execute API for players to remotely control game
}
