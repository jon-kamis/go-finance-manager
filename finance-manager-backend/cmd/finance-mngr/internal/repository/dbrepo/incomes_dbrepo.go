package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"strings"
	"time"
)

func (m *PostgresDBRepo) GetAllUserIncomes(userId int, search string) ([]*models.Income, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "incomes_dbrepo.GetAllUserIncomes"
	fmt.Printf("[ENTER %s]\n", method)

	var query string
	var err error
	var rows *sql.Rows

	if search != "" {
		search = strings.ToLower(search)
		fmt.Printf("[%s] Searching for incomes meeting criteria: %s\n", method, search)
		query = `
		SELECT
			id, user_id, name, type, rate, hours, amount, frequency, tax_percentage, start_dt,
			create_dt, last_update_dt
		FROM incomes
		WHERE
			user_id = $1
			AND
			LOWER(name) like '%' || $2 || '%'`
		rows, err = m.DB.QueryContext(ctx, query, userId, search)
	} else {
		query = `
		SELECT
			id, user_id, name, type, rate, hours, amount, frequency, tax_percentage, start_dt,
			create_dt, last_update_dt
		FROM incomes
		WHERE
			user_id = $1`
		rows, err = m.DB.QueryContext(ctx, query, userId)
	}

	incomes := []*models.Income{}
	recordCount := 0

	if err != nil {
		if err == sql.ErrNoRows {
			return incomes, nil
		} else {
			fmt.Printf("[%s] database call returned with error %s\n", method, err)
			fmt.Printf("[EXIT %s]\n", method)
			return nil, err
		}

	}

	defer rows.Close()

	for rows.Next() {
		var income models.Income
		err := rows.Scan(
			&income.ID,
			&income.UserID,
			&income.Name,
			&income.Type,
			&income.Rate,
			&income.Hours,
			&income.GrossPay,
			&income.Frequency,
			&income.TaxPercentage,
			&income.StartDt,
			&income.CreateDt,
			&income.LastUpdateDt,
		)

		if err != nil {
			fmt.Printf("[%s] error occured when attempting to scan rows into objects\n", method)
			fmt.Printf("[EXIT %s]\n", method)
			return nil, err
		}

		recordCount = recordCount + 1
		incomes = append(incomes, &income)
	}

	fmt.Printf("[%s] retrieved %d records\n", method, recordCount)
	fmt.Printf("[EXIT %s]\n", method)
	return incomes, nil
}

func (m *PostgresDBRepo) GetIncomeByID(id int) (models.Income, error) {
	method := "incomes_dbrepo.GetIncomeByID"
	fmt.Printf("[ENTER %s]\n", method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, user_id, name, type, rate, hours, amount, frequency, tax_percentage, start_dt,
			create_dt, last_update_dt
		FROM incomes
		WHERE 
			id = $1`

	var income models.Income
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&income.ID,
		&income.UserID,
		&income.Name,
		&income.Type,
		&income.Rate,
		&income.Hours,
		&income.GrossPay,
		&income.Frequency,
		&income.TaxPercentage,
		&income.StartDt,
		&income.CreateDt,
		&income.LastUpdateDt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return income, nil
		} else {
			fmt.Printf("[%s] database call returned with error %s\n", method, err)
			fmt.Printf("[EXIT %s]\n", method)
			return income, err
		}

	}

	fmt.Printf("[EXIT %s]\n", method)
	return income, nil
}

func (m *PostgresDBRepo) UpdateIncome(income models.Income) error {
	method := "incomes_dbrepo.UpdateIncome"
	fmt.Printf("[ENTER %s]\n", method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`UPDATE incomes 
		SET
			name = $2,
			type = $3,
			rate = $4,
			hours = $5,
			amount = $6,
			frequency = $7,
			tax_percentage = $8,
			start_dt = $9,
			last_update_dt = $10
		WHERE
			id = $1`

	_, err := m.DB.ExecContext(ctx, stmt,
		income.ID,
		income.Name,
		income.Type,
		income.Rate,
		income.Hours,
		income.GrossPay,
		income.Frequency,
		income.TaxPercentage,
		income.StartDt,
		time.Now(),
	)

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when updating income\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		return err
	}

	fmt.Printf("[EXIT %s]\n", method)
	return nil
}

func (m *PostgresDBRepo) InsertIncome(income models.Income) (int, error) {
	method := "incomes_dbrepo.InsertIncome"
	fmt.Printf("[ENTER %s]\n", method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`INSERT INTO incomes 
			(user_id, name, type, rate, hours, amount, frequency, tax_percentage, start_dt, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt,
		income.UserID,
		income.Name,
		income.Type,
		income.Rate,
		income.Hours,
		income.GrossPay,
		income.Frequency,
		income.TaxPercentage,
		income.StartDt,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		fmt.Printf("[%s] threw error %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return -1, err
	}

	fmt.Println("[EXIT  " + method + "]")
	return id, nil
}

func (m *PostgresDBRepo) DeleteIncomeByID(id int) error {
	method := "incomes_dbrepo.DeleteIncomeByID"
	fmt.Printf("[ENTER %s]\n", method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM incomes
		WHERE 
			id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		fmt.Printf("[%s] database call returned with error %s\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return err
	}

	fmt.Printf("[EXIT %s]\n", method)
	return nil
}
