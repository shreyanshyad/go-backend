package middlewares

import (
	"net/http"

	l "github.com/sirupsen/logrus"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("Request at ", r.Method, " ", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
