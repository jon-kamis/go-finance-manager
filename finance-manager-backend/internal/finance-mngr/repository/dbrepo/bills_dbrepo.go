package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"strings"
	"time"

	"github.com/jon-kamis/klogger"
)

func (m *PostgresDBRepo) GetAllUserBills(userId int, search string) ([]*models.Bill, error) {
	method := "bills_dbrepo.GetAllUserBills"
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
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
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
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

		recordCount = recordCount + 1
		bills = append(bills, &bill)
	}

	klogger.Debug(method, "retrieved %d records\n", recordCount)
	klogger.Exit(method)
	return bills, nil
}

func (m *PostgresDBRepo) GetBillByID(id int) (models.Bill, error) {
	method := "bills_dbrepo.GetBillsByID"
	klogger.Enter(method)

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
			klogger.Exit(method)
			return bill, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return bill, err
		}
	}

	klogger.Exit(method)
	return bill, nil
}

func (m *PostgresDBRepo) UpdateBill(bill models.Bill) error {
	method := "bills_dbrepo.UpdateBill"
	klogger.Enter(method)

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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	klogger.Exit(method)
	return nil
}

func (m *PostgresDBRepo) InsertBill(bill models.Bill) (int, error) {
	method := "bills_dbrepo.InsertBill"
	klogger.Enter(method)

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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return -1, err
	}

	klogger.Exit(method)
	return id, nil
}

func (m *PostgresDBRepo) DeleteBillByID(id int) error {
	method := "bills_dbrepo.DeleteBillByID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM bills
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

func (m *PostgresDBRepo) DeleteBillsByUserID(id int) error {
	method := "bills_dbrepo.DeleteBillsByUserID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM bills
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
