package handlers

import (
	"bank-app/pkg/config"
	"bank-app/pkg/models"
	"bank-app/render"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Customers is the handler for the customers page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "customers.page.tmpl", &models.TemplateData{})
}
