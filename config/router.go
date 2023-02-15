package config

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lucasbravi2019/arquitectura/api/budget"
	"github.com/lucasbravi2019/arquitectura/api/dimension"
	"github.com/lucasbravi2019/arquitectura/api/material"
	"github.com/lucasbravi2019/arquitectura/core"
	"github.com/lucasbravi2019/arquitectura/middleware"
)

var apiRouterInstance *mux.Router

func GetRouter() *mux.Router {
	if apiRouterInstance == nil {
		apiRouterInstance = mux.NewRouter()
	}
	return apiRouterInstance
}

func RegisterRoutes(routes core.Routes) {
	router := GetRouter()
	for _, route := range routes {
		router.
			Path(route.Path).
			HandlerFunc(
				middleware.RequestLoggerMiddleware(
					middleware.DatabaseCheckMiddleware(
						route.HandlerFunc))).
			Methods(route.Method)
	}
}

func StartApi() {
	RegisterRoutes(budget.GetBudgetHandlerInstance().GetBudgetRoutes())
	RegisterRoutes(material.GetMaterialHandlerInstance().GetMaterialRoutes())
	RegisterRoutes(dimension.GetDimensionHandlerInstance().GetDimensionRoutes())

	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(credentials, methods, ttl, origins)(GetRouter())))
}
