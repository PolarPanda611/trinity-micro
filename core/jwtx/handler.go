package jwtx

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PolarPanda611/trinity-micro/core/e"
	"github.com/PolarPanda611/trinity-micro/core/httpx"
	"github.com/dgrijalva/jwt-go"
)

type ContextKey string

const (
	header      string     = "Authorization"
	prefix      string     = "Bearer"
	JwtxContext ContextKey = "jwtx-context"
)

func JwtHandler(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authString := r.Header.Get(header)
			token := strings.Split(authString, " ")
			if len(token) != 2 || token[0] != prefix {
				httpx.HttpResponseErr(r.Context(), w, e.NewError(e.Info, e.ErrUnauthorized, "authorization invalid format"))
				return
			}

			v, err := jwt.Parse(token[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				httpx.HttpResponseErr(r.Context(), w, e.NewError(e.Info, e.ErrUnauthorized, fmt.Sprintf("token parsing error: %v ", err)))
				return
			}
			if !v.Valid {
				httpx.HttpResponseErr(r.Context(), w, e.NewError(e.Info, e.ErrUnauthorized, " invalid token"))
				return
			}
			claim, ok := v.Claims.(jwt.MapClaims)
			if !ok {
				httpx.HttpResponseErr(r.Context(), w, e.NewError(e.Info, e.ErrUnauthorized, " invalid token"))
				return
			}
			r = r.WithContext(
				InjectCtx(r.Context(), claim),
			)
			next.ServeHTTP(w, r)
		})
	}
}

func InjectCtx(ctx context.Context, claim jwt.MapClaims) context.Context {
	return context.WithValue(ctx, JwtxContext, claim)
}

// LoggerFromCtx
// if not exist will panic
func FromCtx(ctx context.Context) jwt.MapClaims {
	claim, ok := ctx.Value(JwtxContext).(jwt.MapClaims)
	if !ok {
		panic("please use logx.InitLogger to init logger ")
	}
	return claim
}
