package player

import (
	"net/http"

	"github.com/LinMAD/TheMazeRunnerServer/maze"
)

// TokenID unique id of player to trace actions
type TokenID string

// HTTPServerAPI used to communicate with players
type HTTPServerAPI struct {
	server   *http.Server
	hostname string
	// Players in this game
	Players map[TokenID]*Player
	// gameMap where players located
	gameMap *maze.Graph
	// mazeMap divided by nodes for player
	mazeMap GameMapData
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

// GameMapData
type GameMapData struct {
	// EncodedMazeNodes for players, contains maze nodes and neighbors
	EncodedMazeNodes []string `json:"maze_nodes"`
}

// response for HTTP
type response struct {
	// Message provided in body
	Message interface{} `json:"message"`
	// Status of HTTP code
	Status int `json:"status"`
}
