package logger

import (
	"fmt"
    "log"
    "net/http"
    "time"
)


func Logger(inner http.Handler, name string) http.Handler {
	fmt.Println("Logger")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}