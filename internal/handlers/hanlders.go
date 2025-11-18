package handlers

import (
	"bank-app/internal/config"
	"bank-app/internal/driver"
	"bank-app/internal/models"
	"bank-app/internal/repository"
	"bank-app/internal/repository/dbrepo"
	"bank-app/render"
	"fmt"
	"net/http"
	"time"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Customers is the handler for the customers page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	customers, err := m.DB.AllCustomers()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("the customer", customers)
	data := make(map[string]any)
	data["customers"] = customers

	render.RenderTemplate(w, "customers.page.tmpl", &models.TemplateData{
		Data: data,
	})
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

	user := models.User{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Username:  r.Form.Get("username"),
		Email:     r.Form.Get("username"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = m.DB.InsertUser(user)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
