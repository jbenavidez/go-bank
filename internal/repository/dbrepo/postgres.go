package dbrepo

import (
	"bank-app/internal/models"
	"context"
	"fmt"
	"time"
)

const dbTimeout = time.Second * 3

func (m *postgresDBRepo) AllCustomers() ([]*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		select
			id, first_name, last_name, email, username
		from
			users
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerList []*models.User

	for rows.Next() {
		var customer models.User
		err := rows.Scan(
			&customer.ID,
			&customer.FirstName,
			&customer.LastName,
			&customer.Email,
			&customer.Username,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		customerList = append(customerList, &customer)
	}
	return customerList, nil
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

func (m *postgresDBRepo) Getuser(userID int) (*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		select
			id, first_name, last_name, email, username
		from
			users
		where id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, userID)
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil

}

func (m *postgresDBRepo) UpdateUser(userID int, userObj models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		update users set first_name=$1, last_name=$2 , email = $3, username = $4, updated_at = $5
		where id = $6
	`
	_, err := m.DB.ExecContext(ctx, stmt,
		userObj.FirstName,
		userObj.LastName,
		userObj.Email,
		userObj.Username,
		userObj.UpdatedAt,
		userID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func (m *postgresDBRepo) DeleteUser(userID int) error {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		Delete From users
		where id = $1
	`
	_, err := m.DB.ExecContext(ctx, stmt,
		userID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
