package middleware

import "net/http"

type middleware func(http.Handler) http.Handler

func Chain(f http.Handler, mf ...middleware) http.Handler {
	if len(mf) == 0 {
		return f
	}

	return mf[0](Chain(f, mf[1:]...))
}
