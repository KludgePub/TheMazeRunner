package player

import (
	"context"
	"net/http"
	"strconv"

	"github.com/LinMAD/TheMazeRunnerServer/maze"
)

const logTag = "-> HTTP API server:"

// TokenID unique id of player to trace actions
type TokenID string

// HTTPServerAPI used to communicate with players
type HTTPServerAPI struct {
	server   *http.Server
	hostname string
	// Players in this game
	Players map[TokenID]*Player
}

// Player general data
type Player struct {
	// ID of player
	ID TokenID `json:"id"`
	// Location of the player
	Location maze.Point `json:"location"`
	// LastMovementPath how player requested to walk
	LastMovementPath []maze.Point `json:"last_movement_path"`
}

// NewPlayerApi to handle players requests
func NewPlayerApi(hostname string) *HTTPServerAPI {
	return &HTTPServerAPI{
		hostname: hostname,
		Players:  make(map[TokenID]*Player),
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
