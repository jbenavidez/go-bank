package config

import (
	"html/template"
	"log"
)

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
}
