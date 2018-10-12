package http

import (
	"io"
	"net/http"

	"github.com/thingful/iotpolicystore/pkg/postgres"
	goji "goji.io"
	"goji.io/pat"
)

// MuxHandlers binds a handler function to the passed in multiplexer. By
// inverting this to set here rather than when creating the server it keeps the
// configuration closer to the handler.
func MuxHandlers(mux *goji.Mux, db *postgres.DB) {
	mux.Handle(pat.Get("/pulse"), healthCheckHandler(db))
}

// healthCheckHandler is a simple http handler that writes `ok` to the
// requester.
func healthCheckHandler(db *postgres.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			http.Error(w, "failed to connect to DB", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "ok")
	})
}
