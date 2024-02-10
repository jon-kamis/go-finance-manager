package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"strings"

	"github.com/jon-kamis/klogger"
)

func (m *PostgresDBRepo) GetAllRoles(search string) ([]*models.Role, error) {
	method := "roles_dbrepo.GetAllRoles"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query string
	var err error
	var rows *sql.Rows

	if search != "" {
		search = strings.ToLower(search)
		query = `
		SELECT
			id, code, create_dt, last_update_dt
		FROM roles
		WHERE
			LOWER(code) like '%' || $1 || '%'`
		rows, err = m.DB.QueryContext(ctx, query, search)
	} else {
		query = `
		SELECT
			id, code, create_dt, last_update_dt
		FROM roles`
		rows, err = m.DB.QueryContext(ctx, query)
	}

	roles := []*models.Role{}
	recordCount := 0

	if err != nil {
		if err == sql.ErrNoRows {
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return roles, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.Code,
			&role.CreateDt,
			&role.LastUpdateDt,
		)

		if err != nil {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

		recordCount = recordCount + 1
		roles = append(roles, &role)
	}

	klogger.Debug(method, "found %d results", recordCount)
	klogger.Exit(method)
	return roles, nil
}

func (m *PostgresDBRepo) GetRoleByCode(code string) (*models.Role, error) {
	method := "roles_dbrepo.GetRoleByCode"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, code, create_dt, last_update_dt
		FROM roles
		WHERE code = $1`

	row := m.DB.QueryRowContext(ctx, query, code)

	var role models.Role
	err := row.Scan(
		&role.ID,
		&role.Code,
		&role.CreateDt,
		&role.LastUpdateDt,
	)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return nil, err
	}

	klogger.Exit(method)
	return &role, nil
}

func (m *PostgresDBRepo) GetRoleById(id string) (*models.Role, error) {
	method := "roles_dbrepo.GetRoleById"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, code, create_dt, last_update_dt
		FROM roles
		WHERE id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var role models.Role
	err := row.Scan(
		&role.ID,
		&role.Code,
		&role.CreateDt,
		&role.LastUpdateDt,
	)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return nil, err
	}

	klogger.Exit(method)
	return &role, nil
}
