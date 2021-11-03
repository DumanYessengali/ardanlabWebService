package mid

import (
	"context"
	"errors"
	"github.com/DumanYessengali/ardanlabWebService/business/auth"
	"github.com/DumanYessengali/ardanlabWebService/foundation/web"
	"log"
	"net/http"
	"strings"
)

var ErrForbidden = web.NewRequestError(
	errors.New("you are not authorized for that action"),
	http.StatusForbidden,
)

func Authenticate(a *auth.Auth) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			authStr := r.Header.Get("Authorization")

			parts := strings.Split(authStr, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: Bearer <token>")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			ctx = context.WithValue(ctx, auth.Key, claims)
			return handler(ctx, w, r)
		}
		return h
	}
	return m
}

func Authorize(log *log.Logger, roles ...string) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, ok := ctx.Value(auth.Key).(auth.Claims)
			if !ok {
				return errors.New("claims missing from context: HasRole called without/before Authentication")
			}
			if !claims.Authorize(roles...) {
				log.Println("mid: authorize: claims: %v exp: %v", claims.Roles, roles)
				return ErrForbidden
			}
			return handler(ctx, w, r)
		}
		return h
	}
	return m
}
