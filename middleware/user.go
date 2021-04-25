package middleware

import (
	"context"
	"net/http"

	"github.com/kotojo/life-manager/models"
)

type userContextKey string

var (
	UserContextKey = userContextKey("user")
)

func User(u *models.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("remember_token")
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			user, err := u.ByRemember(c.Value)
			if err == models.ErrInvalidCredentials {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}
			if err != nil {
				http.Error(w, "Error retrieving user", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func UserFromContext(ctx context.Context) *models.User {
	raw, _ := ctx.Value(UserContextKey).(*models.User)
	return raw
}
