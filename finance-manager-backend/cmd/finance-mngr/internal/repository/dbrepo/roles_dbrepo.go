package dbrepo

import (
	"context"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
)

func (m *PostgresDBRepo) GetUserRoles(id int) ([]*models.UserRole, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "roles_dbrepo.GetUserRoles"
	fmt.Println("[ENTER " + method + "]")

	query := `select id, user_id, role_id, code, create_dt, last_update_dt
		FROM user_roles
		WHERE user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil {
		fmt.Printf("[%s] %s\n", method, "database call returned with error")
		fmt.Println("[EXIT  " + method + "]")
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
