package middleware

import (
	"log"
	"net/http"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		jwtToken := strings.TrimPrefix(authToken, "Bearer ")
		token, err := VerifyTok(jwtToken)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusForbidden)
			_ = token
			return
		}
		next.ServeHTTP(w, r)
	})
}
