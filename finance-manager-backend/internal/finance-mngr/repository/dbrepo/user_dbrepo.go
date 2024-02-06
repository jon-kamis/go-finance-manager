package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"strings"
	"time"

	"github.com/jon-kamis/klogger"
	"golang.org/x/crypto/bcrypt"
)

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	method := "user_dbrepo.GetUserByID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	klogger.Debug(method, "searching for user with id: %d", id)

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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return nil, err
	}

	klogger.Exit(method)
	return &user, nil
}

func (m *PostgresDBRepo) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "user_dbrepo.GetUserByUsername"
	klogger.Enter(method)

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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return nil, err
	}

	klogger.Exit(method)
	return &user, nil
}

func (m *PostgresDBRepo) GetUserByUsernameOrEmail(username string, email string) (*models.User, error) {
	method := "user_dbrepo.GetUserByUsernameOrEmail"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, username, email, first_name, last_name, password,
		create_dt, last_update_dt
		FROM users
		WHERE 
			username = $1
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

		if err == sql.ErrNoRows {
			return &user, nil
		} else {
			klogger.Error(method, constants.UnexpectedSQLError, err)
			return nil, err
		}
	}

	klogger.Exit(method)
	return &user, nil
}

func (m *PostgresDBRepo) InsertUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "user_dbrepo.InsertUser"
	klogger.Enter(method)

	// Encrypt Password
	encryptedPass, bcErr := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if bcErr != nil {
		klogger.ExitError(method, "error occured while encrypting password")
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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return -1, err
	}

	klogger.Exit(method)
	return id, nil
}

func (m *PostgresDBRepo) UpdateUserDetails(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "user_dbrepo.InsertUser"
	klogger.Enter(method)

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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	klogger.Exit(method)
	return nil
}

func (m *PostgresDBRepo) GetAllUsers(search string) ([]*models.User, error) {
	method := "user_dbrepo.GetAllUsers"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query string
	var err error
	var rows *sql.Rows

	if search != "" {
		search = strings.ToLower(search)
		klogger.Debug(method, "searching for user meeting criteria: %s", search)
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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
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
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

		recordCount = recordCount + 1
		users = append(users, &user)
	}

	klogger.Exit(method)
	return users, nil
}

func (m *PostgresDBRepo) DeleteUserByID(id int) error {
	method := "user_dbrepo.DeleteUserByID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM users
		WHERE 
			id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	klogger.Exit(method)
	return nil
}
