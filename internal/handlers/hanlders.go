package handlers

import (
	"bank-app/internal/config"
	"bank-app/internal/driver"
	"bank-app/internal/forms"
	"bank-app/internal/models"
	"bank-app/internal/repository"
	"bank-app/internal/repository/dbrepo"
	"bank-app/render"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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
	data := make(map[string]any)
	data["customers"] = customers

	render.RenderTemplate(w, r, "customers.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AddCustomer is the handler for the Customer add form page
func (m *Repository) AddCustomer(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]any)
	render.RenderTemplate(w, r, "add_customer.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// OpenAccount is the handler for the Open account page
func (m *Repository) OpenAccount(w http.ResponseWriter, r *http.Request) {

	//get customers
	customers, err := m.DB.AllCustomers()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "unable to pull customer list")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	data := make(map[string]any)
	data["customers"] = customers
	render.RenderTemplate(w, r, "add_account.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// CreateAccount is the handler fo creating the account on our db
func (m *Repository) CreateAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Println("something break getting the form")
		m.App.Session.Put(r.Context(), "error", "unable to pull customer list")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	userID, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		fmt.Println(err)
		m.App.Session.Put(r.Context(), "error", "unable to convert userID")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	amount, err := strconv.ParseFloat(r.Form.Get("amount"), 64)
	if err != nil {
		fmt.Println(err)
		m.App.Session.Put(r.Context(), "error", "unable to parse string")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	//set account
	account := models.Account{
		User:        models.User{ID: userID},
		AccountType: r.Form.Get("account_type"),
		Amount:      amount,
	}
	// validate form
	form := forms.New(r.PostForm)
	form.Required("user_id", "account_type", "amount")
	if !form.Valid() {
		data := make(map[string]any)
		data["account"] = account
		render.RenderTemplate(w, r, "add_account.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return

	}
	//create account
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()
	_, err = m.DB.CreateAccount(account)
	if err != nil {
		fmt.Println(err)
		m.App.Session.Put(r.Context(), "error", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	// set flash messagew
	m.App.Session.Put(r.Context(), "flash", "Account was created succefully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CreateCustomer is the handler for creating a customer record when the user submit the form
func (m *Repository) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		fmt.Println("something break getting the form")
		return
	}
	// set user
	user := models.User{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Username:  r.Form.Get("username"),
		Email:     r.Form.Get("email_address"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// validate form
	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email_address", "username")
	form.MinLength("first_name", 1)
	form.IsEmail("email_address")
	form.MinLength("username", 1)
	if !form.Valid() {
		data := make(map[string]any)
		data["user"] = user
		render.RenderTemplate(w, r, "add_customer.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return

	}
	_, err = m.DB.InsertUser(user)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// add flash message before redirect
	m.App.Session.Put(r.Context(), "flash", "user was created succefully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// EditCustomer is the handler for displaying  the edit form
func (m *Repository) EditCustomer(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println(err)
		return
	}

	user, err := m.DB.Getuser(userID)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]any)
	data["customer"] = user
	render.RenderTemplate(w, r, "edit_customers.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ViewCustomer(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println(err)
		return
	}
	// get user
	user, err := m.DB.Getuser(userID)
	if err != nil {
		fmt.Println(err)
		return
	}
	// get accounts
	accounts, err := m.DB.AllAccountsByUserID(userID)
	if err != nil {
		fmt.Println(err)
		return
	}
	//set map
	data := make(map[string]any)
	data["customer"] = user // we could get user info from the accounts | this query is optional
	data["accounts"] = accounts
	render.RenderTemplate(w, r, "view_customers.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// UpdateCustomer is the handler for updating customer re
func (m *Repository) UpdateCustomer(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = r.ParseForm()

	if err != nil {
		fmt.Println("something break getting the form")
		return
	}
	user := models.User{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Username:  r.Form.Get("username"),
		Email:     r.Form.Get("email_address"),
		UpdatedAt: time.Now(),
	}
	// update user
	err = m.DB.UpdateUser(userID, user)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteCustomer is the handler for deleting a customer record
func (m *Repository) DeleteCustomer(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println(err)
		return
	}

	//delete user
	err = m.DB.DeleteUser(userID)
	if err != nil {
		fmt.Println(err)
		// return bad equest
		_ = m.WriteResponse(w, false, "unable to deleted the user", http.StatusBadRequest)
		return
	}
	//set response
	_ = m.WriteResponse(w, false, "user deleted", http.StatusOK)

}

func (m *Repository) WriteResponse(w http.ResponseWriter, errStatus bool, message string, status int) error {
	//set response
	resp := JSONResponse{
		Error:   errStatus,
		Message: message,
	}

	out, err := json.Marshal(resp)

	if err != nil {
		fmt.Println(err)
		return err
	}
	//write header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	//write response
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (m *Repository) Accounts(w http.ResponseWriter, r *http.Request) {

	accounts, err := m.DB.AllAccounts()
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]any)
	data["accounts"] = accounts

	render.RenderTemplate(w, r, "accounts.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
