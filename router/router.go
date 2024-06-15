package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/rohan031/identity/controllers"
)

func Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/identity", controllers.GetIdentity)

	return router
}
