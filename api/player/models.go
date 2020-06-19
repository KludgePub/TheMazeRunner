package player

import (
	"net/http"

	"github.com/LinMAD/TheMazeRunnerServer/manager"
	"github.com/LinMAD/TheMazeRunnerServer/maze"
)

const headerPlayerTokenId = "PlayerID"

// TokenID unique id of player to trace actions
type TokenID string

// HTTPServerAPI used to communicate with players
type HTTPServerAPI struct {
	server   *http.Server
	hostname string
	// Players in this game
	Players    map[TokenID]*Player
	// mazeRawMap detailed information about maze
	mazeRawMap *maze.Map
	// mazeRawGraph structured as graph
	mazeRawGraph *maze.Graph
	// mazeMap divided by nodes for player
	mazeMap GameMapData
	gm *manager.GameManager
}

// Identity data
type Identity struct {
	// Name of player
	Name string `json:"name"`
	// ID assigned to player
	ID TokenID `json:"id"`
}

// Player general data
type Player struct {
	// Identity of the player
	Identity Identity `json:"identity"`
	// Location of the player
	Location maze.Point `json:"location"`
}

// GameMapData
type GameMapData struct {
	// EncodedMazeNodes for players, contains maze nodes and neighbors
	EncodedMazeNodes []string `json:"maze_nodes"`
	// Locations in the maze, start, finish, key
	Locations map[string]maze.Point `json:"locations"`
}

// response for HTTP
type response struct {
	// Message provided in body
	Message interface{} `json:"message"`
	// Status of HTTP code
	Status int `json:"status"`
}
