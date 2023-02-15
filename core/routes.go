package core

import "net/http"

type Route struct {
	Path        string
	HandlerFunc func(w http.ResponseWriter, r *http.Request)
	Method      string
}

type Routes []Route
