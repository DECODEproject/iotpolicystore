package http

import (
	"io"
	"net/http"

	goji "goji.io"
	"goji.io/pat"
)

// MuxHandlers binds a handler function to the passed in multiplexer. By
// inverting this to set here rather than when creating the server it keeps the
// configuration closer to the handler.
func MuxHandlers(mux *goji.Mux) {
	mux.HandleFunc(pat.Get("/pulse"), healthCheckHandler)
}

// healthCheckHandler is a simple http handler that writes `ok` to the
// requester.
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "ok")
}
