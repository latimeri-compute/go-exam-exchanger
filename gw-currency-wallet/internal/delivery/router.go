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
		// Authorization: Bearer JWT_TOKEN
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(utils.JWTAuthenticator())

			r.Get("/balance", h.GetBalance)
			r.Post("/deposit", h.TopUpBalance)
			r.Post("/withdraw", h.WithdrawFromBalance)

			r.Get("/exchange/rates", h.GetExchangeRates)
			r.Post("/exchange", h.ExchangeFunds)
		})
	})

	return r
}
