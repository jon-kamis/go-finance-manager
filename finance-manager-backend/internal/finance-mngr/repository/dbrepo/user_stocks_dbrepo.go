package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/constants"
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

	var stmt string
	var err error
	var id int

	//Set default type
	if s.Type == "" {
		s.Type = constants.UserStockTypeOwn
	}

	if !s.ExpirationDt.Time.IsZero() {
		stmt =
			`INSERT INTO user_stocks 
			(user_id, ticker, quantity, type, effective_dt, expiration_dt, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5, $6, $7, $8) returning id`

		err = m.DB.QueryRowContext(ctx, stmt,
			s.UserId,
			s.Ticker,
			s.Quantity,
			s.Type,
			s.EffectiveDt,
			s.ExpirationDt.Time,
			time.Now(),
			time.Now(),
		).Scan(&id)
	} else {
		stmt =
			`INSERT INTO user_stocks 
				(user_id, ticker, quantity, type, effective_dt, create_dt, last_update_dt)
			values 
				($1, $2, $3, $4, $5, $6, $7) returning id`

		err = m.DB.QueryRowContext(ctx, stmt,
			s.UserId,
			s.Ticker,
			s.Quantity,
			s.Type,
			s.EffectiveDt,
			time.Now(),
			time.Now(),
		).Scan(&id)
	}

	if err != nil {
		fmlogger.ExitError(method, "error occured when inserting new user stock", err)
		return -1, err
	}

	fmlogger.Exit(method)
	return id, nil
}

// Function GetAllUserStocks returns all user stocks with names matching search if it is included and where t is after effective dt and before expiration date if it exists
func (m *PostgresDBRepo) GetAllUserStocks(userId int, stockType string, search string, t time.Time) ([]*models.UserStock, error) {
	method := "stocks_dbrepo.GetAllUserStocks"
	fmlogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query string
	var err error
	var rows *sql.Rows

	//Set to default stock type
	if stockType == "" {
		stockType = constants.UserStockTypeOwn
	}

	if search != "" {
		search = strings.ToLower(search)

		query = `
		SELECT
			id, user_id, ticker, quantity, effective_dt, expiration_dt,
			create_dt, last_update_dt
		FROM user_stocks
		WHERE
			user_id = $1
			AND
			type=$2
			AND
			effective_dt <= $3
			AND
			(expiration_dt IS NULL OR expiration_dt >= $3)
			AND
			LOWER(ticker) like '%' || $4 || '%'`
		rows, err = m.DB.QueryContext(ctx, query, userId, stockType, t, search)
	} else {
		query = `
		SELECT
			id, user_id, ticker, quantity, effective_dt, expiration_dt,
			create_dt, last_update_dt
		FROM user_stocks
		WHERE
			user_id = $1
			AND
			type = $2
			AND
			effective_dt <= $3
			AND
			(expiration_dt IS NULL OR expiration_dt >= $3)`
		rows, err = m.DB.QueryContext(ctx, query, userId, stockType, t)
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
			&u.EffectiveDt,
			&u.ExpirationDt,
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

// Function GetAllUserStocksByDateRange returns all user stocks with names matching search if it is included and where the userStock was active during any part of the date range
func (m *PostgresDBRepo) GetAllUserStocksByDateRange(userId int, search string, ts time.Time, te time.Time) ([]*models.UserStock, error) {
	method := "stocks_dbrepo.GetAllUserStocks"
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
			id, user_id, ticker, quantity, effective_dt, expiration_dt,
			create_dt, last_update_dt
		FROM user_stocks
		WHERE
			user_id = $1
			AND
			effective_dt <= $3
			AND
			(expiration_dt IS NULL OR expiration_dt >= $2)
			AND
			LOWER(ticker) like '%' || $4 || '%'`
		rows, err = m.DB.QueryContext(ctx, query, userId, ts, te, search)
	} else {
		query = `
		SELECT
			id, user_id, ticker, quantity, effective_dt, expiration_dt,
			create_dt, last_update_dt
		FROM user_stocks
		WHERE
			user_id = $1
			AND
			effective_dt <= $3
			AND
			(expiration_dt IS NULL OR expiration_dt >= $2)`
		rows, err = m.DB.QueryContext(ctx, query, userId, ts, te)
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
			&u.EffectiveDt,
			&u.ExpirationDt,
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
