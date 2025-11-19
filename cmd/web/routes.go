package main

import (
	"bank-app/internal/config"
	"bank-app/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Get("/", handlers.Repo.Home)

	mux.Get("/customers/add", handlers.Repo.AddCustomer)
	mux.Post("/customers/create", handlers.Repo.CreateCustomer)

	mux.Get("/customers/{id}", handlers.Repo.EditCustomer)
	mux.Post("/customers/{id}", handlers.Repo.UpdateCustomer)
	mux.Delete("/customers/delete/{id}", handlers.Repo.DeleteCustomer)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
