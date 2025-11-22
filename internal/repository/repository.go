package repository

import "bank-app/internal/models"

type DatabaseRepo interface {
	AllCustomers() ([]*models.User, error)
	InsertUser(user models.User) (int, error)
	Getuser(userID int) (*models.User, error)
	UpdateUser(userID int, userObj models.User) error
	DeleteUser(userID int) error
	CreateAccount(account models.Account) (int, error)
	AllAccounts() ([]*models.Account, error)
	AllAccountsByUserID(userId int) ([]*models.Account, error)
}
