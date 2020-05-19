package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/LinMAD/TheMazeRunnerServer/maze"
	"github.com/LinMAD/TheMazeRunnerServer/validator"
)

func init() {
	fmt.Printf("\n%s\n", "-> Initilizing the maze server...")
	rand.Seed(time.Now().UnixNano())
}

func main() {
	fmt.Println("-> Generating new maze...")

	// Size of maze map
	r, c := 3,3
	m := maze.NewMaze(r, c)
	m.Generate()

	fmt.Println("-> Maze ready...")
	fmt.Println("- Bytes map:\n", m.Container)
	fmt.Println("- Visual map: ")
	fmt.Println(maze.PrintMaze(m))

	// TODO Debug mode
	fmt.Println("\nDebug maze rows output")
	g := maze.DispatchToGraph(m)
	for id, n := range g.Nodes {
		fmt.Printf("\nNode: %d with prop: %s => neighbors:\n%s", id, string(n.Entity), maze.PrintGraphNode(n))
	}

	fmt.Println("\nDebug solved maze")
	solved := validator.SolveMaze(*m)
	fmt.Println(maze.PrintMaze(solved.SolvedMap))
	fmt.Printf("\nPath can be solved (%v) from Start to Key:\n%v", validator.IsPathPossible(solved.ToKey, g), solved.ToKey)
	fmt.Printf("\nPath can be solved (%v) from Key to Exit:\n%v", validator.IsPathPossible(solved.ToExit, g),solved.ToExit)

	// TODO Execute API for players to remotely control game
}
