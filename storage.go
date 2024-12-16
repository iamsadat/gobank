package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Storage is an interface that defines the methods that must be implemented by a storage layer.
type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

// PostgresStore is a storage layer that uses Postgres as the underlying database.
type PostgresStore struct {
	db *sql.DB
}

// // DeleteAccount implements Storage.
// func (s *PostgresStore) DeleteAccount(int) error {
// 	panic("unimplemented")
// }

// // GetAccountByID implements Storage.
// func (s *PostgresStore) GetAccountByID(int) (*Account, error) {
// 	panic("unimplemented")
// }

// // UpdateAccount implements Storage.
// func (s *PostgresStore) UpdateAccount(*Account) error {
// 	panic("unimplemented")
// }

// NewPostgresStore creates a new PostgresStore
func NewPostgresStore() (*PostgresStore, error) {
	// connStr is the connection string used to connect to the Postgres database
	connStr := "postgresql://neondb_owner:6NPyaLzZwH9s@ep-long-unit-a1aos54v.ap-southeast-1.aws.neon.tech/neondb?sslmode=require"
	// sql.Open opens a new database connection using the provided driver and source name
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping verifies a connection to the database is still alive, establishing a connection if necessary
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// the function returns a pointer to a PostgresStore instance and nil error
	return &PostgresStore{
		db: db,
	}, nil
}

// init function initializes the PostgresStore instance by calling the createAccountTable function which creates a table in the database
func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

// createAccountTable creates a table in the database if it does not exist
func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		number INTEGER,
		balance NUMERIC,
		created_at TIMESTAMP
	)`

	_, err := s.db.Exec(query)
	return err
}

// CreateAccount is a query that is used to create a new account in the database
func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
	INSERT INTO account 
	(first_name, last_name, number, balance, created_at)
	VALUES 
	($1, $2, $3, $4, $5)`
	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}
	fmt.Printf("resp: %v\n", resp)
	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM account WHERE ID = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE number = $1", number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with number [%d] not found", number)
}

// GetAccountByID is a query that is used to search for a specific account in the database by the account id
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	row, err := s.db.Query("SELECT * FROM account WHERE ID = $1", id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		return scanIntoAccount(row)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

// GetAccounts is a query that is used to fetch all the accounts in the database
func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// scanIntoAccount is a helper function that is used to scan the rows returned from the database into an account struct
func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)

	return account, err
}
