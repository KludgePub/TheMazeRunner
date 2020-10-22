package player

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/LinMAD/TheMazeRunner/generator"
	"github.com/LinMAD/TheMazeRunner/manager"
	"github.com/LinMAD/TheMazeRunner/maze"
	"github.com/LinMAD/TheMazeRunner/validator"
)

// jsonResponse helper to wrap responses to json format
func (api *HTTPServerAPI) jsonResponse(w http.ResponseWriter, msg interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;utf-8")

	b, err := json.Marshal(response{Message: msg, Status: status})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message": "Service Unavailable", "status": "503"}`))

		log.Printf("%s unable to marshal response to json error: %v", logTag, err)
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(b)
}

func (api *HTTPServerAPI) extractPlayer(id TokenID) (p *Player, err error) {
	if p, ok := api.Players[id]; ok {
		return p, nil
	}

	log.Printf("%s unable to find player by given token: %s", logTag, id)

	return nil, fmt.Errorf("player not found")
}

//
// Controller handlers
//

// handlerHomeDoc will provide with list of API calls
func (api *HTTPServerAPI) handlerHomeDoc(w http.ResponseWriter, r *http.Request) {
	// TODO Return actions lists for routes, how to get world, how to make a move etc.

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 {
		api.jsonResponse(w, "Hello, "+parts[2], http.StatusOK)
		return
	}

	api.jsonResponse(w, "Hello, world", http.StatusOK)
}

// handlerPlayerRegister to current game session
func (api *HTTPServerAPI) handlerPlayerRegister(w http.ResponseWriter, r *http.Request) {
	var p Identity
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		log.Printf("%s unable to get body from accept path route: %v", logTag, err)
		api.jsonResponse(w, "Invalid data in body", http.StatusBadRequest)
		return
	}

	id, idErr := generator.CreateUUID()
	if idErr != nil {
		log.Printf("%s unable to create new uuid: %v", logTag, idErr)
		api.jsonResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	p.ID = TokenID(id)
	api.Players[p.ID] = &Player{
		Identity: p,
		Location: api.mazeRawMap.Entrance,
		gameOver: false,
	}

	api.jsonResponse(w, p, http.StatusOK)
}

// handlerPlayerStats returns information about player // TODO Scoring of players?
func (api *HTTPServerAPI) handlerPlayerStats(w http.ResponseWriter, r *http.Request) {
	pid := r.Header.Get(headerPlayerTokenId)
	p, pErr := api.extractPlayer(TokenID(pid))
	if pErr != nil {
		api.jsonResponse(w, "Player ID not found, register your self first", http.StatusBadRequest)
		return
	}

	api.jsonResponse(w, p, http.StatusOK)
}

// handlerPlayerAcceptPath collect requested movement path from player
func (api *HTTPServerAPI) handlerPlayerAcceptPath(w http.ResponseWriter, r *http.Request) {
	pid := r.Header.Get(headerPlayerTokenId)
	p, pErr := api.extractPlayer(TokenID(pid))
	if pErr != nil {
		api.jsonResponse(w, "Player ID not found, register your self first", http.StatusBadRequest)
		return
	}

	var path []maze.Point
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&path); err != nil {
		log.Printf("%s unable to get body from accept path route: %v", logTag, err)
		api.jsonResponse(w, "Invalid data in body", http.StatusBadRequest)
		return
	}

	path = append([]maze.Point{p.Location}, path...)
	log.Printf("%s Player %s given path: %v", logTag, pid, path)
	possiblePath := validator.GetPossiblePath(path, api.mazeRawGraph)
	log.Printf("%s Player %s possible path: %v", logTag, pid, possiblePath)
	p.Location = possiblePath[len(possiblePath)-1]
	log.Printf("%s Player %s new location: %v", logTag, pid, p.Location)

	api.gm.AddMovement(manager.PlayerInfo{
		Name:            p.Identity.Name,
		ID:              string(p.Identity.ID),
		CurrentLocation: p.Location,
		MovementPath:    possiblePath,
	})

	api.jsonResponse(w, "Movement path accepted", http.StatusAccepted)
}

// handlerPlayerInteraction player can interact in location
func (api *HTTPServerAPI) handlerPlayerInteraction(w http.ResponseWriter, r *http.Request) {
	pid := r.Header.Get(headerPlayerTokenId)
	p, pErr := api.extractPlayer(TokenID(pid))
	if pErr != nil {
		api.jsonResponse(w, "Player ID not found, register your self first", http.StatusBadRequest)
		return
	}

	// Check if player in location where interaction possible
	if p.Location != api.mazeRawMap.Key && p.Location != api.mazeRawMap.Exit {
		api.jsonResponse(w, "Nothing to do here", http.StatusOK)
		return
	}

	// Respond by action method
	// Get item id
	if r.Method == http.MethodGet && p.Location == api.mazeRawMap.Key {
		api.jsonResponse(w, Item{ID: api.mazeRawMap.KeyCode}, http.StatusOK)
		return
	}

	// Check if all conditions passed for exit
	if r.Method == http.MethodPost && p.Location == api.mazeRawMap.Exit {
		var item Item
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&item); err != nil {
			log.Printf("%s unable to get body from accept path route: %v", logTag, err)
			api.jsonResponse(w, "Invalid data in body", http.StatusBadRequest)
			return
		}
		if api.mazeRawMap.KeyCode != item.ID {
			api.jsonResponse(w, "Invalid item id or expired", http.StatusBadRequest)
			return
		}

		p.gameOver = true

		api.jsonResponse(w, "Congratulations you did it!", http.StatusOK)
		return
	}

	api.jsonResponse(w, "Are you sure about what are you doing?", http.StatusBadRequest)
}

// handlerWorldMazeData give to player maze data to assemble: current node, right node, bottom node => 0 0, 0 1, 1 0  (x,y)
func (api *HTTPServerAPI) handlerWorldMazeData(w http.ResponseWriter, r *http.Request) {
	api.jsonResponse(w, api.mazeMap, http.StatusOK)
}
