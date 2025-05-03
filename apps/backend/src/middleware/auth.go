package middleware

import (
	"context"
	"net/http"

	"github.com/okzmo/kyob/db"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "failed to find token", http.StatusUnauthorized)
			return
		}

		user, err := db.Query.VerifyToken(r.Context(), token.Value)
		if err != nil {
			http.Error(w, "failed to verify token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
