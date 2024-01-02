package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"time"
)

func (m *PostgresDBRepo) GetAllUserRoles(id int) ([]*models.UserRole, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "userRoles_dbrepo.GetUserRoles"
	fmt.Println("[ENTER " + method + "]")

	query := `select id, user_id, role_id, code, create_dt, last_update_dt
		FROM user_roles
		WHERE user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil {
		fmt.Printf("[%s] database call returned with error\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		return nil, err
	}
	defer rows.Close()

	var userRoles []*models.UserRole
	recordCount := 0

	for rows.Next() {
		var userRole models.UserRole
		err := rows.Scan(
			&userRole.ID,
			&userRole.UserId,
			&userRole.RoleId,
			&userRole.Code,
			&userRole.CreateDt,
			&userRole.LastUpdateDt,
		)

		if err != nil {
			fmt.Printf("[%s] %s\n", method, "returned with error")
			fmt.Println("[EXIT  " + method + "]")
			return nil, err
		}
		recordCount = recordCount + 1
		userRoles = append(userRoles, &userRole)
	}

	fmt.Printf("[%s] loaded %d roles for user with id %d\n", method, recordCount, id)

	fmt.Println("[EXIT  " + method + "]")
	return userRoles, nil
}

func (m *PostgresDBRepo) GetUserRoleById(id int) (models.UserRole, error) {
	method := "userRoles_dbRepo.GetUserRoleById"
	fmlogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, user_id, role_id, code, create_dt, last_update_dt
		FROM user_roles
		WHERE id = $1`

	var userRole models.UserRole
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&userRole.ID,
		&userRole.UserId,
		&userRole.RoleId,
		&userRole.Code,
		&userRole.CreateDt,
		&userRole.LastUpdateDt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmlogger.Info(method, constants.EntityNotFoundError)
			fmlogger.Exit(method)
			return userRole, nil
		} else {
			fmlogger.ExitError(method, constants.UnexpectedSQLError, err)
			return userRole, err
		}

	}

	fmlogger.Exit(method)
	return userRole, nil
}

func (m *PostgresDBRepo) InsertUserRole(userRole models.UserRole) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "userRoles_dbrepo.InsertUserRole"
	fmt.Printf("[ENTER %s]", method)

	stmt :=
		`INSERT INTO user_roles 
			(user_id, role_id, code, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5) returning id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt,
		userRole.UserId,
		userRole.RoleId,
		userRole.Code,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		fmt.Printf("[%s] error when attempting to save new user role\n", method)
		fmt.Printf("[EXIT %s]", method)
		return -1, err
	}

	fmt.Printf("[EXIT %s]", method)
	return id, nil
}

func (m *PostgresDBRepo) DeleteUserRolesByUserID(id int) error {
	method := "userRoles_dbRepo.DeleteUserRolesByUserID"
	fmlogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM user_roles
		WHERE 
			user_id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		fmlogger.ExitError(method, "database call returned with error", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}

func (m *PostgresDBRepo) DeleteUserRoleByID(id int) error {
	method := "userRoles_dbRepo.DeleteUserRoleByID"
	fmlogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM user_roles
		WHERE 
			id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		fmlogger.ExitError(method, "database call returned with error", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}
