package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "user_dbrepo.GetUserByID"
	fmt.Println("[ENTER " + method + "]")
	fmt.Printf("[%s] %s %d\n", method, "searching for user with id ", id)

	query := `select id, username, email, first_name, last_name, password,
		create_dt, last_update_dt
		FROM users
		WHERE id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreateDt,
		&user.LastUpdateDt,
	)

	if err != nil {
		fmt.Printf("[%s] %s\n", method, "returned with error")
		fmt.Println("[EXIT  " + method + "]")
		return nil, err
	}

	fmt.Println("[EXIT  " + method + "]")
	return &user, nil
}

func (m *PostgresDBRepo) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "user_dbrepo.GetUserByUsername"
	fmt.Println("[ENTER " + method + "]")

	query := `select id, username, email, first_name, last_name, password,
		create_dt, last_update_dt
		FROM users
		WHERE username =$1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreateDt,
		&user.LastUpdateDt,
	)

	if err != nil {
		fmt.Printf("[%s] %s\n", method, "returned with error")
		fmt.Println("[EXIT  " + method + "]")
		return nil, err
	}

	fmt.Println("[EXIT  " + method + "]")
	return &user, nil
}

func (m *PostgresDBRepo) GetUserByUsernameOrEmail(username string, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "user_dbrepo.GetUserByUsernameOrEmail"
	fmt.Println("[ENTER " + method + "]")

	query := `select id, username, email, first_name, last_name, password,
		create_dt, last_update_dt
		FROM users
		WHERE 
			username =$1
			OR email = $2`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, username, email)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreateDt,
		&user.LastUpdateDt,
	)

	if err != nil {
		fmt.Printf("[%s] %s\n", method, "returned with error")
		fmt.Println("[EXIT  " + method + "]")
		return nil, err
	}

	fmt.Println("[EXIT  " + method + "]")
	return &user, nil
}

func (m *PostgresDBRepo) InsertUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "user_dbrepo.InsertUser"
	fmt.Println("[ENTER " + method + "]")

	// Encrypt Password
	encryptedPass, bcErr := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if bcErr != nil {
		fmt.Printf("[%s] %s\n", method, "error occured while encrypting password")
		fmt.Println("[EXIT  " + method + "]")
		return -1, bcErr
	}

	stmt :=
		`INSERT INTO users 
			(username, email, first_name, last_name,
			password, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5, $6, $7) returning id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		string(encryptedPass),
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		fmt.Printf("[%s] %s\n", method, "returned with error")
		fmt.Println("[EXIT  " + method + "]")
		return -1, err
	}

	fmt.Println("[EXIT  " + method + "]")
	return id, nil
}

func (m *PostgresDBRepo) UpdateUserDetails(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "user_dbrepo.InsertUser"
	fmt.Println("[ENTER " + method + "]")

	stmt :=
		`UPDATE users 
		SET
			email = $2,
			first_name = $3,
			last_name = $4,
			last_update_dt = $5)
		WHERE
			id = $1`

	err := m.DB.QueryRowContext(ctx, stmt,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		time.Now(),
	).Scan()

	if err != nil {
		fmt.Printf("[%s] %s\n", method, "returned with error")
		fmt.Println("[EXIT  " + method + "]")
		return err
	}

	fmt.Println("[EXIT  " + method + "]")
	return nil
}

func (m *PostgresDBRepo) GetAllUsers(search string) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "user_dbrepo.GetAllUsers"
	fmt.Printf("[ENTER %s]\n", method)

	var query string
	var err error
	var rows *sql.Rows

	if search != "" {
		search = strings.ToLower(search)
		fmt.Printf("[%s] Searching for user meeting criteria: %s\n", method, search)
		query = `
		SELECT
			id, username, email, first_name, last_name, password,
			create_dt, last_update_dt
		FROM users
		WHERE 
			LOWER(username) like '%' || $1 || '%'
			OR LOWER(first_name) like '%' || $1 || '%'
			OR LOWER(last_name) like '%' || $1 || '%'
			OR LOWER(email) like '%' || $1 || '%'`
		rows, err = m.DB.QueryContext(ctx, query, search)
	} else {
		query = `
		SELECT
			id, username, email, first_name, last_name, password,
			create_dt, last_update_dt
		FROM users`
		rows, err = m.DB.QueryContext(ctx, query)
	}

	var users []*models.User
	recordCount := 0

	if err != nil {
		fmt.Printf("[%s] database call returned with error %s\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.CreateDt,
			&user.LastUpdateDt,
		)

		if err != nil {
			fmt.Printf("[%s] error occured when attempting to scan rows into objects\n", method)
			fmt.Printf("[EXIT %s]\n", method)
			return nil, err
		}

		recordCount = recordCount + 1
		users = append(users, &user)
	}

	fmt.Printf("[EXIT %s]\n", method)
	return users, nil
}
