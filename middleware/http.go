package middleware

import (
	"context"
	"net/http"
)

var (
	HttpContextKey = "http"
)

type HttpContextValues struct {
	W *http.ResponseWriter
	R *http.Request
}

func Http() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpContextValues := HttpContextValues{
				W: &w,
				R: r,
			}
			ctx := context.WithValue(r.Context(), HttpContextKey, httpContextValues)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
