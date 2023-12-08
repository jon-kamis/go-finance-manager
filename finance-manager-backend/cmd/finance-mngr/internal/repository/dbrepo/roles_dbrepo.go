package dbrepo

import (
	"context"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
)

func (m *PostgresDBRepo) GetRoleByCode(code string) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "roles_dbrepo.AddUserRole"
	fmt.Println("[ENTER " + method + "]")

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
		fmt.Printf("[%s] %s\n", method, "returned with error")
		fmt.Println("[EXIT  " + method + "]")
		return nil, err
	}

	fmt.Printf("[%s] loaded roles successfully\n", method)

	fmt.Println("[EXIT  " + method + "]")
	return &role, nil
}
