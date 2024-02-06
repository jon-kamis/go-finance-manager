package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"time"

	"github.com/jon-kamis/klogger"
)

func (m *PostgresDBRepo) GetAllUserRoles(id int) ([]*models.UserRole, error) {
	method := "userRoles_dbrepo.GetUserRoles"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, user_id, role_id, code, create_dt, last_update_dt
		FROM user_roles
		WHERE user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
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
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}
		recordCount = recordCount + 1
		userRoles = append(userRoles, &userRole)
	}

	klogger.Debug(method, "loaded %d roles for user with id %d", recordCount, id)

	klogger.Exit(method)
	return userRoles, nil
}

func (m *PostgresDBRepo) GetUserRoleByID(id int) (models.UserRole, error) {
	method := "userRoles_dbRepo.GetUserRoleById"
	klogger.Enter(method)

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
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return userRole, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return userRole, err
		}

	}

	klogger.Exit(method)
	return userRole, nil
}

func (m *PostgresDBRepo) GetUserRoleByRoleIDAndUserID(rId int, uId int) (models.UserRole, error) {
	method := "userRoles_dbRepo.GetUserRoleByRoleId"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, user_id, role_id, code, create_dt, last_update_dt
		FROM user_roles
		WHERE 
			role_id = $1
			AND user_id = $2`

	var userRole models.UserRole
	row := m.DB.QueryRowContext(ctx, query, rId, uId)

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
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return userRole, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return userRole, err
		}

	}

	klogger.Exit(method)
	return userRole, nil
}

func (m *PostgresDBRepo) InsertUserRole(userRole models.UserRole) (int, error) {
	method := "userRoles_dbrepo.InsertUserRole"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return -1, err
	}

	klogger.Enter(method)
	return id, nil
}

func (m *PostgresDBRepo) DeleteUserRolesByUserID(id int) error {
	method := "userRoles_dbRepo.DeleteUserRolesByUserID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM user_roles
		WHERE 
			user_id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	klogger.Exit(method)
	return nil
}

func (m *PostgresDBRepo) DeleteUserRoleByID(id int) error {
	method := "userRoles_dbRepo.DeleteUserRoleByID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM user_roles
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

func (m *PostgresDBRepo) DeleteUserRoleByRoleID(id int) error {
	method := "userRoles_dbRepo.DeleteUserRoleByRoleID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM user_roles
		WHERE 
			role_id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	klogger.Exit(method)
	return nil
}
