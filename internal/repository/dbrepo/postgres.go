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

func (m *postgresDBRepo) CreateAccount(account models.Account) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `
		insert into accounts
			(account_type, amount, user_id, created_at, updated_at)
		values	
			($1,$2,$3,$4,$5) 
		returning id
	`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		account.AccountType,
		account.Amount,
		account.User.ID,
		account.CreatedAt,
		account.CreatedAt,
	).Scan(&newID)

	if err != nil {
		fmt.Println("err creating account", err)
		return 0, err
	}
	return newID, nil

}

func (m *postgresDBRepo) AllAccounts() ([]*models.Account, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		select
			accounts.id, accounts.account_type,  accounts.account_status , accounts.amount, accounts.created_at, accounts.updated_at, users.id, users.first_name, users.last_name 
		from
			accounts
		inner join
			users ON accounts.user_id = users.id

	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*models.Account

	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.ID,
			&account.AccountType,
			&account.AccountStatus,
			&account.Amount,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.User.ID,
			&account.User.FirstName,
			&account.User.LastName,
		)
		if err != nil {
			fmt.Println("error getting accounts", err)
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

func (m *postgresDBRepo) AllAccountsByUserID(userId int) ([]*models.Account, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		select
			accounts.id, accounts.account_type, accounts.account_status , accounts.amount, accounts.created_at, accounts.updated_at, users.id, users.first_name, users.last_name 
		from
			accounts
		inner join
			users ON accounts.user_id = users.id
		where 
			 accounts.user_id = $1
	`
	rows, err := m.DB.QueryContext(ctx, query,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*models.Account

	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.ID,
			&account.AccountType,
			&account.AccountStatus,
			&account.Amount,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.User.ID,
			&account.User.FirstName,
			&account.User.LastName,
		)
		if err != nil {
			fmt.Println("error getting accounts", err)
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

func (m *postgresDBRepo) GetAccount(accountId int) (*models.Account, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		select
			id, account_type, amount, created_at, updated_at, account_status
		from
			accounts
		where id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, accountId)
	var account models.Account
	err := row.Scan(
		&account.ID,
		&account.AccountType,
		&account.Amount,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.AccountStatus,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &account, nil

}

func (m *postgresDBRepo) UpdateAccount(accountID int, account models.Account) error {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		update accounts set account_type=$1, account_status=$2
		where id = $3
	`
	_, err := m.DB.ExecContext(ctx, stmt,
		account.AccountType,
		account.AccountStatus,
		accountID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
