package player

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

// handlerPlayerMazeData give to player maze data to assemble: current node, right node, bottom node => 0 0, 0 1, 1 0  (x,y)
func (api *HTTPServerAPI) handlerPlayerMazeData(w http.ResponseWriter, r *http.Request) {
	api.jsonResponse(w, api.mazeMap.EncodedMazeNodes, http.StatusOK)
}

// TODO Register new player
// TODO Collect requested movement path
// TODO Interaction handling with maze object like with "K" key to collect and handle race condition between players
// TODO Scoring of players?
