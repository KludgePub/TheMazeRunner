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

	// Size of maze map
	r, c := 2, 2
	m := maze.NewMaze(r, c)
	m.Generate()

	fmt.Println("-> Maze ready...")
	fmt.Println("- Bytes map:\n", m.Container)
	fmt.Println("- Visual map: ")
	fmt.Println(maze.PrintMaze(m))
	for id, n := range maze.DispatchToGraph(m).Nodes {
		fmt.Printf("\nNode: %d with prop: %s => neighbors:\n%s", id, string(n.Entity), maze.PrintGraphNode(n))
	}

	// TODO Execute API for players to remotely control game
}
