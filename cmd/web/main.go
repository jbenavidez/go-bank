package main

import (
	"bank-app/internal/config"
	"bank-app/internal/driver"
	"bank-app/internal/handlers"
	"bank-app/internal/helpers"
	"bank-app/render"
	"fmt"
	"log"
	"net/http"
	"os"
)

const port = ":8080"

var app config.AppConfig

// logger
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

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
	//

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 user=postgres password=postgres dbname=bank sslmode=disable timezone=UTC connect_timeout=5")
	if err != nil {
		log.Fatal("Cannot connect to database!")
	}
	log.Println("Connected to database!")
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
