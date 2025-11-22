package models

import "time"

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Account struct {
	ID            int     `json:"id"`
	AccountType   string  `json:"account_type"`
	Amount        float64 `json:"amount"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	User          User
	AccountStatus string
}
