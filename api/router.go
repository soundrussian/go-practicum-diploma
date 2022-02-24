package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/soundrussian/go-practicum-diploma/api/handlers"
	customMiddleware "github.com/soundrussian/go-practicum-diploma/api/middleware"
)

func routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(customMiddleware.LogRequest)

	r.Post("/api/user/register", handlers.HandleRegister)

	return r
}
