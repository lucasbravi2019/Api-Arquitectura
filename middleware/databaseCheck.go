package middleware

import (
	"log"
	"net/http"

	"github.com/lucasbravi2019/arquitectura/core"
)

func DatabaseCheckMiddleware(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := core.CheckDatabaseHealth()
		if err != nil {
			log.Println(err)
			return
		}
		f(w, r)
	}
}
