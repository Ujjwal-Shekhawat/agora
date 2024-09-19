package routes

import (
	"gateway_service/api/controllers"
	"net/http"
)

func RegisterControllers(mux *http.ServeMux, controllers ...controllers.Controller) {
	for _, c := range controllers {
		c.RegisterRoutes(mux)
	}
}
