package config

import (
	"fmt"
	"net/http"
)

func (cfg *Config) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method,
			r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (cfg *Config) PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defered functions are always called in case of panic when GO unwinds the stack

		defer func() {
			// err created by recover() can be of any type, depending what was passed to panic()
			// mostly err, string
			if err := recover(); err != nil {
				// Connection close - header acts as a trigger to autmatically close gos http connection
				w.Header().Set("Connection", "close")
            cfg.Hlp.ServerError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
