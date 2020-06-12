package player

import (
	"fmt"
	"log"
	"net/http"
)

// bootRouter boots router
func (api *HTTPServerAPI) bootRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", api.router)

	return mux
}

// router for the requests
func (api *HTTPServerAPI) router(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s Incoming request (ContentLength: %d) - URL: %s", logTag, r.ContentLength, r.URL.Path)

	var pid string
	path := r.URL.Path

	if path != "/player/new" && path != "/" {
		pid = r.Header.Get(headerPlayerTokenId)
		if pid == "" || len(pid) < 16 {
			api.jsonResponse(
				w,
				fmt.Sprintf("Not found %s in request headers or id is invalid", headerPlayerTokenId),
				http.StatusBadRequest,
			)
			return
		} else if _, extErr := api.extractPlayer(TokenID(pid)); extErr != nil {
			api.jsonResponse(
				w,
				fmt.Sprintf("Not found %s in request headers or id is invalid", headerPlayerTokenId),
				http.StatusBadRequest,
			)

			return
		}
	}

	log.Printf("%s Handling request from (%s) to %s\n", logTag, pid, r.URL.Path)

	// TODO Loopback requesters if they don't have player token id
	switch r.Method {
	case http.MethodPost:
		switch {
		case path == "/player/new":
			api.handlerPlayerRegister(w, r)
		case path == "/player/move":
			api.handlerPlayerAcceptPath(w, r)
		default:
			api.jsonResponse(w, "not found", http.StatusNotFound)
		}
	case http.MethodGet:
		switch {
		case path == "/player/stats":
			api.handlerPlayerStats(w, r)
		case path == "/player/world":
			api.handlerPlayerMazeData(w, r)
		case path == "/":
			api.handlerHomeDoc(w, r)
		default:
			api.jsonResponse(w, "not found", http.StatusNotFound)
		}
	default:
		api.jsonResponse(w, "not found", http.StatusNotFound)
	}
}
