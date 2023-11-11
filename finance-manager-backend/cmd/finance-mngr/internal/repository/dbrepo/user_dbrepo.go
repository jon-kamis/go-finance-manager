package dbrepo

import (
	"context"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
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

	method := "user_dbrepo.GetUserByUsername"
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

func (m *PostgresDBRepo) InsertUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "user_dbrepo.InsertUser"
	fmt.Println("[ENTER " + method + "]")

	// Encrypt Password
	encryptedPass, bcErr := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if bcErr != nil {
		fmt.Printf("[%s] %s\n", method, "error occured while encrypting password")
		fmt.Println("[EXIT  " + method + "]")
		return bcErr
	}

	stmt :=
		`INSERT INTO users 
			(username, email, first_name, last_name,
			password, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		string(encryptedPass),
		time.Now(),
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
