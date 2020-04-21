package main

import (
	"fmt"
	"github.com/LinMAD/TheMazeRunnerServer/maze"
	"math/rand"
	"time"
)

func main()  {
	fmt.Printf("\n%s\n\n", "The maze runner")

	rand.Seed(time.Now().UnixNano())

	m := maze.NewMaze(10, 10)
	m.Generate()

	fmt.Println(m)
}
