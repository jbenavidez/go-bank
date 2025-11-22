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
	mux.Use(SessionLoadAndSave) //middleware to use session

	mux.Get("/", handlers.Repo.Home)
	// customers
	mux.Get("/customers/add", handlers.Repo.AddCustomer)
	mux.Post("/customers/create", handlers.Repo.CreateCustomer)
	mux.Get("/customers/{id}", handlers.Repo.EditCustomer)
	mux.Post("/customers/{id}", handlers.Repo.UpdateCustomer)
	mux.Delete("/customers/delete/{id}", handlers.Repo.DeleteCustomer)
	mux.Get("/customers/view/{id}", handlers.Repo.ViewCustomer)
	//accounts
	mux.Get("/accounts", handlers.Repo.Accounts)
	mux.Get("/accounts/create", handlers.Repo.OpenAccount)
	mux.Post("/accounts/create", handlers.Repo.CreateAccount)
	mux.Get("/accounts/view/{id}", handlers.Repo.ViewAccount)
	mux.Get("/accounts/edit/{id}", handlers.Repo.EditAccount)
	mux.Post("/accounts/edit/{id}", handlers.Repo.UpdateAccount)
	//
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
