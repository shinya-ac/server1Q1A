package auth0

import (
	"context"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type JWTMiddlewareKey struct{}
type JWTKey struct{}

func WithJWTMiddleware(m *jwtmiddleware.JWTMiddleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), JWTMiddlewareKey{}, m)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UseJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtm, ok := r.Context().Value(JWTMiddlewareKey{}).(*jwtmiddleware.JWTMiddleware)
		if !ok || jwtm == nil {
			logging.Logger.Error("JWT middleware not found")
			http.Error(w, "JWT middleware not found", http.StatusInternalServerError)
			return
		}

		err := jwtm.CheckJWT(w, r)
		if err != nil {
			logging.Logger.Error("JWT verification failed", "error", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		val := r.Context().Value(jwtm.Options.UserProperty)
		if val == nil {
			logging.Logger.Error("UserProperty not found in context")
		} else {
			logging.Logger.Info("UserProperty found in context", "val", val)
			if token, ok := val.(*jwt.Token); ok {
				logging.Logger.Info("JWT token successfully cast")
				ctx := context.WithValue(r.Context(), JWTKey{}, token)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else {
				logging.Logger.Error("Failed to cast UserProperty to *jwt.Token")
			}
		}

		next.ServeHTTP(w, r)
	})
}

func GetJWT(ctx context.Context) *jwt.Token {
	rawJWT, ok := ctx.Value(JWTKey{}).(*jwt.Token)
	if !ok {
		return nil
	}
	return rawJWT
}
