package main

import (
	"bank-app/internal/config"
	"bank-app/internal/driver"
	"bank-app/internal/handlers"
	"bank-app/internal/helpers"
	"bank-app/internal/models"
	"bank-app/render"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const port = ":8080"

var app config.AppConfig
var sessionManager *scs.SessionManager

func main() {
	//  data that is going to session
	//gob.Register(models.User{})

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Starting application on port %s", port))

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	gob.Register(models.User{})
	// set session
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction
	//add session to app
	app.Session = sessionManager
	// setdb
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 user=postgres password=postgres dbname=bank sslmode=disable timezone=UTC connect_timeout=5")
	if err != nil {
		log.Fatal("unable to connect to database!")
	}
	log.Println("Connected to database!")

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
