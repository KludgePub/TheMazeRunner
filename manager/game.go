package manager

import (
	"encoding/json"
	"sync"

	"github.com/LinMAD/TheMazeRunnerServer/maze"
)

// GameManager
type GameManager struct {
	m *maze.Map
	// pl players list with id as key
	pl   map[string]PlayerInfo
	lock sync.RWMutex
}

// PlayerInfo
type PlayerInfo struct {
	// Name of player
	Name string `json:"name"`
	// ID unique key
	ID   string `json:"id"`
	// CurrentLocation of the player
	CurrentLocation maze.Point `json:"current"`
	// MovementPath next destination
	MovementPath []maze.Point `json:"path"`
	// TODO For the score can be counted submits ratio as modifier
}

// NewGameManager will handle and communicate throw APIs
func NewGameManager(mazeMap *maze.Map) *GameManager {
	return &GameManager{pl: make(map[string]PlayerInfo, 0)}
}

// AddMovement of player
func (gm *GameManager) AddMovement(info PlayerInfo) {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	if _, ok := gm.pl[info.ID]; !ok {
		gm.pl[info.ID] = info
	} else {
		gm.pl[info.ID] = info
	}
}

func (gm *GameManager) GetPlayerMovements() (moves []byte) {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	moves = make([]byte, 0)

	for id, path := range gm.pl {
		b, err := json.Marshal(path)
		if err != nil {
			return
		}

		delete(gm.pl, id)

		moves = b
		return
	}

	return
}
