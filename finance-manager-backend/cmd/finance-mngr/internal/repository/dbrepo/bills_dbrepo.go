package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"strings"
	"time"
)

func (m *PostgresDBRepo) GetAllUserBills(userId int, search string) ([]*models.Bill, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	method := "bills_dbrepo.GetAllUserBills"
	fmlogger.Enter(method)

	var query string
	var err error
	var rows *sql.Rows

	if search != "" {
		search = strings.ToLower(search)

		query = `
		SELECT
			id, user_id, name, amount,
			create_dt, last_update_dt
		FROM bills
		WHERE
			user_id = $1
			AND
			LOWER(name) like '%' || $2 || '%'`
		rows, err = m.DB.QueryContext(ctx, query, userId, search)
	} else {
		query = `
		SELECT
			id, user_id, name, amount,
			create_dt, last_update_dt
		FROM bills
		WHERE
			user_id = $1`
		rows, err = m.DB.QueryContext(ctx, query, userId)
	}

	bills := []*models.Bill{}
	recordCount := 0

	if err != nil {
		if err == sql.ErrNoRows {
			return bills, nil
		} else {
			fmlogger.ExitError(method, "database call returned with error", err)
			return nil, err
		}

	}

	defer rows.Close()

	for rows.Next() {
		var bill models.Bill
		err := rows.Scan(
			&bill.ID,
			&bill.UserID,
			&bill.Name,
			&bill.Amount,
			&bill.CreateDt,
			&bill.LastUpdateDt,
		)

		if err != nil {
			fmlogger.ExitError(method, "error occured when attempting to scan db result into rows", err)
			return nil, err
		}

		recordCount = recordCount + 1
		bills = append(bills, &bill)
	}

	fmt.Printf("[%s] retrieved %d records\n", method, recordCount)
	fmlogger.Exit(method)
	return bills, nil
}

func (m *PostgresDBRepo) GetBillByID(id int) (models.Bill, error) {
	method := "bills_dbrepo.GetBillsByID"
	fmlogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, user_id, name, amount,
			create_dt, last_update_dt
		FROM bills
		WHERE 
			id = $1`

	var bill models.Bill
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&bill.ID,
		&bill.UserID,
		&bill.Name,
		&bill.Amount,
		&bill.CreateDt,
		&bill.LastUpdateDt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmlogger.Exit(method)
			return bill, nil
		} else {
			fmlogger.ExitError(method, "database call returned with error", err)
			return bill, err
		}
	}

	fmlogger.Exit(method)
	return bill, nil
}

func (m *PostgresDBRepo) UpdateBill(bill models.Bill) error {
	method := "bills_dbrepo.UpdateBill"
	fmlogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`UPDATE bills 
		SET
			name = $2,
			amount = $3,
			last_update_dt = $4
		WHERE
			id = $1`

	_, err := m.DB.ExecContext(ctx, stmt,
		bill.ID,
		bill.Name,
		bill.Amount,
		time.Now(),
	)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when updating bill", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}

func (m *PostgresDBRepo) InsertBill(bill models.Bill) (int, error) {
	method := "bills_dbrepo.InsertBill"
	fmlogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`INSERT INTO bills 
			(user_id, name, amount, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5) returning id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt,
		bill.UserID,
		bill.Name,
		bill.Amount,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		fmlogger.ExitError(method, "error occured when inserting new bill", err)
		return -1, err
	}

	fmlogger.Exit(method)
	return id, nil
}

func (m *PostgresDBRepo) DeleteBillByID(id int) error {
	method := "bills_dbrepo.DeleteBillByID"
	fmlogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM bills
		WHERE 
			id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		fmlogger.ExitError(method, "error occured when deleting bill", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}

func (m *PostgresDBRepo) DeleteBillsByUserID(id int) error {
	method := "bills_dbrepo.DeleteBillsByUserID"
	fmlogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM bills
		WHERE 
			user_id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		fmlogger.ExitError(method, "error occured when deleting bills", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}
