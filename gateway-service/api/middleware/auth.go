package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const AuthUserString = contextKey("AuthUserString")

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		jwtToken := strings.TrimPrefix(authToken, "Bearer ")
		token, err := VerifyTok(jwtToken)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		subject, err := token.Claims.GetSubject()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), AuthUserString, subject)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
