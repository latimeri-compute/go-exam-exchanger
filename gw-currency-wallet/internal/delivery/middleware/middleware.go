package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
	"github.com/pascaldekloe/jwt"
)

type ContextID uint

var ModuleName string = "github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/"

func JWTAuthenticator(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			claims, err := jwt.HMACCheckHeader(r, secret)
			if err != nil {
				utils.ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			if !claims.Valid(time.Now()) ||
				claims.Issuer != ModuleName || !claims.AcceptAudience(ModuleName) {
				utils.ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			userID, err := strconv.ParseInt(claims.Subject, 10, 64)
			if err != nil {
				utils.InternalErrorResponse(w)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "user", ContextID(userID)))

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
