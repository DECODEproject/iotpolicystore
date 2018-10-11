package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	// RequestIDHeader is a constant defining the request header we read/write ids
	// to
	RequestIDHeader = "X-Request-ID"

	// RequestCtxKey is a key to be used in a context, under which the incoming
	// request ID is stored.
	RequestCtxKey = contextKey("requestID")
)

// RequestIDMiddleware is a net.http middleware that adds a UUID request ID to
// incoming requests. If the request already contains a request ID supplied by
// the remote user, this is used else we generate a new one. The request id is
// added to the context that is passed down to handlers.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(RequestIDHeader)
		if rid == "" {
			rid = uuid.New().String()
		}
		w.Header().Set(RequestIDHeader, rid)
		ctx := context.WithValue(r.Context(), RequestCtxKey, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
