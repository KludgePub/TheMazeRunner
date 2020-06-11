package player

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type response struct {
	// Message provided in body
	Message string `json:"message"`
	// Status of HTTP code
	Status  int    `json:"status"`
}

// jsonResponse helper to wrap responses to json format
func (api *HTTPServerAPI) jsonResponse(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json;utf-8")

	b, err := json.Marshal(response{Message: msg, Status: status})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message": "Service Unavailable", "status": "503"}`))

		log.Printf("%s unable to marshal response to JSON error: %v", logTag, err)
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(b)
}

//
// Controller handlers
//

// handlerHome will provide with list of API calls
func (api *HTTPServerAPI) handlerHome(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 {
		api.jsonResponse(w, "Hello, "+parts[2], http.StatusOK)
		return
	}

	api.jsonResponse(w, "Hello, world", http.StatusOK)
}

// TODO Register new player
// TODO Collect requested movement path
// TODO Interaction handling with maze object like with "K" key to collect and handle race condition between players
// TODO Scoring of players?

