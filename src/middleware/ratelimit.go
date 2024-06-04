package middleware

import (
	"fs/src/config"
	"fs/src/server/response"
	"github.com/zytekaron/gorl"
	"net/http"
)

// RateLimit implements IP-based rate limiting on an endpoint.
//
// Depends on DataReader.
func RateLimit(cfg *config.Config) Middleware {
	tokens := cfg.Server.Tokens
	rlCfg := cfg.Server.RateLimit

	rl := gorl.New(rlCfg.Limit, rlCfg.Burst, rlCfg.Refill)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Context().Value("middleware.data").(*Data)

			// tokens `admin` and `no_ratelimit` bypass ratelimiting
			if header.Token == tokens.NoRatelimit || header.Token == tokens.Admin {
				next.ServeHTTP(w, r)
				return
			}

			b := rl.Get(header.IP)
			if !b.CanDraw(1) {
				response.WriteError(w, http.StatusTooManyRequests)
				return
			}

			b.Draw(1)
			next.ServeHTTP(w, r)
		})
	}
}
