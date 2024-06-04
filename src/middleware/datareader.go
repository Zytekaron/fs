package middleware

import (
	"context"
	"fs/src/server/response"
	"log"
	"net"
	"net/http"
)

type Data struct {
	IP    string
	Port  string
	Token string
}

// DataReader attaches a struct for data values to the request's context.
//
//	middleware.data -> *Data
func DataReader() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := &Data{}
			ctx := context.WithValue(r.Context(), "middleware.data", h)
			r = r.WithContext(ctx)

			var err error
			h.IP, h.Port, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Println("error parsing ip:")
				response.WriteError(w, 500)
				return
			}

			h.Token = r.Header.Get("X-API-Key")
			if h.Token == "" {
				h.Token = r.URL.Query().Get("api_key")
			}

			next.ServeHTTP(w, r)
		})
	}
}
