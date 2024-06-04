package middleware

import (
	"fs/src/server/response"
	"github.com/zytekaron/gorl"
	"net/http"
	"slices"
	"time"
)

// RequireToken requires that the user specifies a permitted tokens.
//
// Depends on DataReader.
func RequireToken(valid []string) Middleware {
	rl := gorl.New(10, 50, time.Minute)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Context().Value("middleware.data").(*Data)
			b := rl.Get(header.IP)

			if !b.CanDraw(1) {
				response.WriteError(w, 429)
				return
			}

			if len(header.Token) == 0 {
				b.Draw(1)
				response.WriteError(w, 403)
				return
			}
			if slices.Contains(valid, header.Token) {
				b.Draw(1)
				response.WriteError(w, 401)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
