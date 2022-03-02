package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	customMiddleware "github.com/soundrussian/go-practicum-diploma/api/middleware"
)

func (api *API) routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(customMiddleware.LogRequest)

	r.Post("/api/user/register", api.HandleRegister)
	r.Post("/api/user/login", api.HandleLogin)

	return r
}
