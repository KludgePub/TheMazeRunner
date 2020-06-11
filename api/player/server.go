package player

import (
	"context"
	"net/http"
	"strconv"

	"github.com/LinMAD/TheMazeRunnerServer/maze"
)

const logTag = "-> HTTP API server:"

// NewPlayerApi to handle players requests
func NewPlayerApi(gameMap *maze.Graph, hostname string) *HTTPServerAPI {
	egm := make([]string, len(gameMap.Nodes))
	i := 0
	for _, n := range gameMap.Nodes {
		egm[i] = maze.PrintGraphNode(n, false)
		i++
	}

	return &HTTPServerAPI{
		hostname: hostname,
		Players:  make(map[TokenID]*Player),
		gameMap:  gameMap,
		mazeMap:  GameMapData{egm},
	}
}

// Start starts the player API server
func (api *HTTPServerAPI) Start(port int) error {
	p := strconv.Itoa(port)

	s := &http.Server{
		Addr:    ":" + p,
		Handler: api.bootRouter(),
	}

	api.server = s

	return s.ListenAndServe()
}

// Shutdown performs graceful closing of players API server
func (api *HTTPServerAPI) Shutdown() error {
	return api.server.Shutdown(context.Background())
}
