package repository

import "bank-app/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertUser(user models.User) (int, error)
}
