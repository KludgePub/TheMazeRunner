package player

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/LinMAD/TheMazeRunnerServer/generator"
	"github.com/LinMAD/TheMazeRunnerServer/manager"
	"github.com/LinMAD/TheMazeRunnerServer/maze"
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
		api.jsonResponse(w, "Internal server error", http.StatusInternalServerError)
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

// handlerPlayerMazeData give to player maze data to assemble: current node, right node, bottom node => 0 0, 0 1, 1 0  (x,y)
func (api *HTTPServerAPI) handlerPlayerMazeData(w http.ResponseWriter, r *http.Request) {
	api.jsonResponse(w, api.mazeMap.EncodedMazeNodes, http.StatusOK)
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
		api.jsonResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	api.gm.AddMovement(manager.PlayerInfo{
		Name:            p.Identity.Name,
		ID:              string(p.Identity.ID),
		CurrentLocation: p.Location,
		MovementPath:    path,
	})

	api.jsonResponse(w, "Movement path accepted", http.StatusAccepted)
}

// TODO Interaction handling with maze object like with "K" key to collect and handle race condition between players
