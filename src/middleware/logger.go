package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Logger keeps track of basic statistics about requests and recovers from panics.
func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					fmt.Println("recovered from panic:", err)
					fmt.Println(debug.Stack())
				}
			}()

			start := time.Now()

			lw := &loggingResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(lw, r)

			end := time.Now()

			fmt.Printf("%s :: %21s :: %d %s :: %s %s :: %s\n",
				end.Format("2006-01-02 15:04:05.99"),
				r.RemoteAddr,
				lw.statusCode,
				http.StatusText(lw.statusCode),
				r.Method,
				r.URL.EscapedPath(),
				end.Sub(start),
			)
		})
	}
}
