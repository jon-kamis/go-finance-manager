package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"fmt"
	"strings"
)

func (m *PostgresDBRepo) GetAllRoles(search string) ([]*models.Role, error) {
	method := "roles_dbrepo.GetAllRoles"
	fmlogger.Enter(method)

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
			//No data
			return roles, nil
		} else {
			fmlogger.ExitError(method, "database call returned with unexpected error", err)
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
			fmlogger.ExitError(method, "exception thrown when scanning sql rows", err)
			return nil, err
		}

		recordCount = recordCount + 1
		roles = append(roles, &role)
	}

	fmt.Printf("[%s] found %d results\n", method, recordCount)
	fmlogger.Exit(method)
	return roles, nil
}

func (m *PostgresDBRepo) GetRoleByCode(code string) (*models.Role, error) {
	method := "roles_dbrepo.GetRoleByCode"
	fmlogger.Enter(method)

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
		fmlogger.ExitError(method, "exited with error", err)
		return nil, err
	}

	fmlogger.Info(method, "loaded roles successfully")
	fmt.Printf("[%s] \n", method)

	fmlogger.Exit(method)
	return &role, nil
}

func (m *PostgresDBRepo) GetRoleById(id string) (*models.Role, error) {
	method := "roles_dbrepo.GetRoleById"
	fmlogger.Enter(method)

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
		fmlogger.ExitError(method, "exited with error", err)
		return nil, err
	}

	fmlogger.Info(method, "loaded roles successfully")
	fmt.Printf("[%s] \n", method)

	fmlogger.Exit(method)
	return &role, nil
}
