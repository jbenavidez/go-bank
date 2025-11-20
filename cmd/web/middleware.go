package main

import (
	"net/http"
)

func SessionLoadAndSave(next http.Handler) http.Handler {
	return sessionManager.LoadAndSave(next)
}
