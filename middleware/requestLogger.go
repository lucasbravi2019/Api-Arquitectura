package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func RequestLoggerMiddleware(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonRequest *json.RawMessage = &json.RawMessage{}

		err := json.NewDecoder(r.Body).Decode(jsonRequest)
		if err != nil && err.Error() != "EOF" {
			log.Println(err.Error())
		}
		r.Body.Close()

		body, err := jsonRequest.MarshalJSON()
		if err != nil {
			log.Println(err)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		log.Printf("URL Request: %s: %s%s\n", r.Method, r.Host, r.URL.Path)
		stringBody := string(body)
		if stringBody != "" {
			log.Println(stringBody)
		}
		fn(w, r)
	}
}
