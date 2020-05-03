package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/LinMAD/TheMazeRunnerServer/maze"
)

func main()  {
	fmt.Printf("\n%s\n\n", "The maze runner")

	rand.Seed(time.Now().UnixNano())

	m := maze.NewMaze(3, 3)
	m.Generate()

	fmt.Println(m)
}
