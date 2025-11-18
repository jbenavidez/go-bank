package dbrepo

import (
	"bank-app/internal/models"
	"context"
	"fmt"
	"time"
)

const dbTimeout = time.Second * 3

func (m *postgresDBRepo) AllUsers() bool {
	return true

}

func (m *postgresDBRepo) InsertUser(user models.User) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `
		insert into Users
			(first_name, last_name, email, username)
		values	
			($1,$2,$3,$4) 
		returning id
	`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Username,
	).Scan(&newID)

	if err != nil {
		fmt.Println("err creating user", err)
		return 0, err
	}
	return newID, nil

}
