package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
	"github.com/pascaldekloe/jwt"
)

type ContextID uint

var ModuleName string = "github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/"

// проверяет полученный JWT в хедере и добавляет его в контекст
func JWTAuthenticator(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			claims, err := jwt.HMACCheckHeader(r, secret)
			if err != nil {
				utils.UnauthorizedResponse(w, "unauthorized")
				return
			}

			if !claims.Valid(time.Now()) ||
				claims.Issuer != ModuleName || !claims.AcceptAudience(ModuleName) {
				utils.UnauthorizedResponse(w, "unauthorized")
				return
			}

			userID, err := strconv.ParseInt(claims.Subject, 10, 64)
			if err != nil {
				utils.InternalErrorResponse(w)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "userId", ContextID(userID)))

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}

// достаёт ID пользователя из контекста и ищет его в бд
// в случае, если пользователя нет в системе -- статус 401
func RetrieveUserFromDB(userModel storages.UserModelInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			id, ok := r.Context().Value("userId").(ContextID)
			if !ok {
				utils.InternalErrorResponse(w)
			}

			var user storages.User
			user.ID = uint(id)
			err := userModel.FindUser(&user)
			if err != nil {
				if errors.Is(err, storages.ErrRecordNotFound) {
					utils.UnauthorizedResponse(w, "Unauthorized")
					return
				} else {
					utils.InternalErrorResponse(w)
					return
				}
			}

			r = r.WithContext(
				context.WithValue(r.Context(), "user", user),
			)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}

func IssueNewJWT(source []byte, userID int) ([]byte, error) {
	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(int64(userID), 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(time.Hour * 24))
	claims.Issuer = ModuleName
	claims.Audiences = []string{ModuleName}

	jwtBytes, err := claims.HMACSign(jwt.HS256, source)
	if err != nil {
		return nil, err
	}
	return jwtBytes, nil
}
