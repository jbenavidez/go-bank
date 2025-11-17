package handlers

import (
	"bank-app/pkg/config"
	"bank-app/pkg/models"
	"bank-app/render"
	"fmt"
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

// AddCustomer is the handler for the Customer add form page
func (m *Repository) AddCustomer(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "add_customer.page.tmpl", &models.TemplateData{})
}

// CreateCustomer is the handler for creating a customer record when the user submit the form
func (m *Repository) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		fmt.Println("something break getting the form")
		return
	}
	fmt.Println("the username", r.Form.Get("first_name"))
}
