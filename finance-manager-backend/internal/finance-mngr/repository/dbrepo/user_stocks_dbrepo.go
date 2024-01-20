package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"fmt"
	"strings"
	"time"
)

func (m *PostgresDBRepo) InsertUserStock(s models.UserStock) (int, error) {
	method := "stocks_dbrepo.InsertUserStock"
	fmlogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`INSERT INTO user_stocks 
			(user_id, ticker, quantity, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5) returning id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt,
		s.UserId,
		s.Ticker,
		s.Quantity,
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

func (m *PostgresDBRepo) GetAllUserStocks(userId int, search string) ([]*models.UserStock, error) {
	method := "bills_dbrepo.GetAllUserStocks"
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
			id, user_id, ticker, quantity,
			create_dt, last_update_dt
		FROM user_stocks
		WHERE
			user_id = $1
			AND
			LOWER(ticker) like '%' || $2 || '%'`
		rows, err = m.DB.QueryContext(ctx, query, userId, search)
	} else {
		query = `
		SELECT
			id, user_id, ticker, quantity,
			create_dt, last_update_dt
		FROM user_stocks
		WHERE
			user_id = $1`
		rows, err = m.DB.QueryContext(ctx, query, userId)
	}

	usl := []*models.UserStock{}
	recordCount := 0

	if err != nil {
		if err == sql.ErrNoRows {
			return usl, nil
		} else {
			fmlogger.ExitError(method, "database call returned with error", err)
			return nil, err
		}

	}

	defer rows.Close()

	for rows.Next() {
		var u models.UserStock
		err := rows.Scan(
			&u.ID,
			&u.UserId,
			&u.Ticker,
			&u.Quantity,
			&u.CreateDt,
			&u.LastUpdateDt,
		)

		if err != nil {
			fmlogger.ExitError(method, "error occured when attempting to scan db result into rows", err)
			return nil, err
		}

		recordCount = recordCount + 1
		usl = append(usl, &u)
	}

	fmt.Printf("[%s] retrieved %d records\n", method, recordCount)
	fmlogger.Exit(method)
	return usl, nil
}
