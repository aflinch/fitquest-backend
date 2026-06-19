package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	jwtlib "fitquest-backend/jwt"
)

type errorResponse struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func writeGraphQLError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{
		Errors: []struct {
			Message string `json:"message"`
		}{{Message: message}},
	})
}

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"user"}

type UserCtx struct {
	UserID   string
	Username string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")

			claims, err := jwtlib.ParseToken(tokenStr)
			if err != nil {
				writeGraphQLError(w, http.StatusForbidden, "Invalid token")
				return
			}

			user := &UserCtx{
				UserID:   claims.UserID,
				Username: claims.Username,
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *UserCtx {
	raw, _ := ctx.Value(userCtxKey).(*UserCtx)
	return raw
}
