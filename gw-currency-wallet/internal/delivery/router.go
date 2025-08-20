package delivery

import (
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
)

func Router(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.AllowContentType("application/json"))
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.LoginUser)

		// 	Заголовки:
		// _Authorization: Bearer JWT_TOKEN_
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(utils.JWTAuthenticator())

			r.Get("/balance", GetBalance)
			r.Post("/deposit", TopUpBalance)
			r.Post("/withdraw", WithdrawFromBalance)

			r.Get("/exchange/rates", GetExchangeRates)
			r.Post("/exchange", ExchangeFunds)
		})
	})

	return r
}
