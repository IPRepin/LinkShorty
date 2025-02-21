package middleware

import (
	"LinkShorty/configs"
	"LinkShorty/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
	if err != nil {
		return
	}
}

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			writeUnauthorized(w)
			return
		}
		if authorization == "" {
			next.ServeHTTP(w, r)
		}
		token := strings.TrimPrefix(authorization, "Bearer ")
		isValid, data := jwt.NewJWT(config.Auth.Secret).ParseToken(token)
		if !isValid {
			writeUnauthorized(w)
		}
		r.Context()
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
