package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	customMiddleware "github.com/soundrussian/go-practicum-diploma/api/middleware"
	auth "github.com/soundrussian/go-practicum-diploma/auth/v1"
)

func (api *API) routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Use(customMiddleware.LogRequest)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))

		r.Post("/api/user/register", api.HandleRegister)
		r.Post("/api/user/login", api.HandleLogin)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(customMiddleware.CurrentUser)

		r.Get("/api/user/balance", api.HandleBalance)
		r.Get("/api/user/balance/withdrawals", api.HandleWithdrawals)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))

			r.Post("/api/user/balance/withdraw", api.HandleWithdraw)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("text/plain"))

			r.Post("/api/user/orders", api.HandleOrder)
		})
	})

	return r
}
