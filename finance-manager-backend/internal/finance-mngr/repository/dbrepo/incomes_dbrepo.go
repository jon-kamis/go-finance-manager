package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"fmt"
	"strings"
	"time"

	"github.com/jon-kamis/klogger"
)

func (m *PostgresDBRepo) GetAllUserIncomes(userId int, search string) ([]*models.Income, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "incomes_dbrepo.GetAllUserIncomes"
	klogger.Enter(method)

	var query string
	var err error
	var rows *sql.Rows

	if search != "" {
		search = strings.ToLower(search)
		klogger.Debug(method, "searching for incomes meeting criteria: %s", search)
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
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
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
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

		recordCount = recordCount + 1
		incomes = append(incomes, &income)
	}

	klogger.Debug(method, "retrieved %d records", recordCount)
	klogger.Exit(method)
	return incomes, nil
}

func (m *PostgresDBRepo) GetIncomeByID(id int) (models.Income, error) {
	method := "incomes_dbrepo.GetIncomeByID"
	klogger.Enter(method)

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
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return income, err
		}

	}

	klogger.Exit(method)
	return income, nil
}

func (m *PostgresDBRepo) UpdateIncome(income models.Income) error {
	method := "incomes_dbrepo.UpdateIncome"
	klogger.Enter(method)

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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	klogger.Exit(method)
	return nil
}

func (m *PostgresDBRepo) InsertIncome(income models.Income) (int, error) {
	method := "incomes_dbrepo.InsertIncome"
	klogger.Enter(method)

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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return -1, err
	}

	fmt.Println("[EXIT  " + method + "]")
	return id, nil
}

func (m *PostgresDBRepo) DeleteIncomeByID(id int) error {
	method := "incomes_dbrepo.DeleteIncomeByID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM incomes
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

func (m *PostgresDBRepo) DeleteIncomesByUserID(id int) error {
	method := "incomes_dbrepo.DeleteIncomesByUserID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM incomes
		WHERE 
			user_id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
	}

	klogger.Exit(method)
	return nil
}
