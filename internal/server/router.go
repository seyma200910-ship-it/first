package server

import (
	"net/http"

	"first/internal/users"
)

type Handlers struct {
	UserHandler *users.Handler
}

func NewRouter(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	// users
	mux.HandleFunc("POST /users", h.UserHandler.CreateUser)

	// health
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}
