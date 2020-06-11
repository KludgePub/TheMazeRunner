package player

import (
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
	pId := "TODO ADD PLAYER TOKEN ID" // TODO Get player token from request headers

	log.Printf(
		"%s Incoming request from %s to URL: %s\n ContentLength: %d",
		logTag,
		r.URL.Path,
		pId,
		r.ContentLength,
	)

	path := r.URL.Path

	switch r.Method {
	case http.MethodGet:
		switch {
		case path == "/":
			api.handlerHome(w, r)
		default:
			api.jsonResponse(w, "not found", http.StatusNotFound)
		}
	default:
		api.jsonResponse(w, "not found", http.StatusNotFound)
	}
}
