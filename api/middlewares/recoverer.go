package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/thedevsaddam/renderer"
)

// Recoverer recover from panic and log to sentry
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {

				if err, ok := rvr.(error); ok {
					log.Println("Panic", err)
				}

				fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
				debug.PrintStack()

				renderer.New().JSON(w, http.StatusInternalServerError, renderer.M{
					"message": "Internal error",
					"error":   rvr,
				})
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
