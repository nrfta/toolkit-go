package dataloaders

import (
	"context"
	"net/http"
)

// contextKey is unexported to avoid collisions in the request context.
type contextKey struct{}

var loadersKey = contextKey{}

// Middleware attaches dataloaders to the request context.
// The provided newLoaders function is invoked for each request and its
// result stored under a private context key.
func Middleware[L any](newLoaders func(*http.Request) L) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loaders := newLoaders(r)
			r = r.WithContext(context.WithValue(r.Context(), loadersKey, loaders))
			next.ServeHTTP(w, r)
		})
	}
}

// For retrieves the dataloaders value from the context. It panics if the value
// is missing or of the wrong type.
func For[L any](ctx context.Context) L {
	return ctx.Value(loadersKey).(L)
}
