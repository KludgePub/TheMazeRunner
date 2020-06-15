package player

import (
	"context"
	"net/http"
	"strconv"

	"github.com/LinMAD/TheMazeRunnerServer/manager"
	"github.com/LinMAD/TheMazeRunnerServer/maze"
)

const logTag = "-> HTTP API server:"

// NewPlayerApi to handle players requests
func NewPlayerApi(gm *manager.GameManager, mazeMap *maze.Map, mazeGraph *maze.Graph, hostname string) *HTTPServerAPI {
	egm := make([]string, mazeMap.Size)

	i := 0
	for _, n := range mazeGraph.Nodes {
		egm[i] = maze.PrintGraphNode(n, false)
		i++
	}

	return &HTTPServerAPI{
		hostname:     hostname,
		Players:      make(map[TokenID]*Player),
		mazeRawMap:   mazeMap,
		mazeRawGraph: mazeGraph,
		mazeMap:      GameMapData{egm},
		gm: gm,
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
