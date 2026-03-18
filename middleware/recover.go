package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"gitlab.com/durfy/durpify/handlers"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg := "Caught panic: %v, Stack trace: %s"
				log.Printf(msg, err, string(debug.Stack()))

				resp := handlers.NewFailureResponse(
					"Unternal Server Error",
					http.StatusInternalServerError,
					[]string{err.(string)},
				)
				resp.SendReponse(w, r)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
