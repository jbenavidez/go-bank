package main

import (
	"bank-app/pkg/config"
	"bank-app/pkg/handlers"
	"bank-app/render"
	"fmt"
	"log"
	"net/http"
)

const port = ":8080"

var app config.AppConfig

func main() {

	app.InProduction = false
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Staring application on port %s", port))

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
