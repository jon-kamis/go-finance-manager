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

func (m *PostgresDBRepo) InsertUserStock(s models.UserStock) (int, error) {
	method := "stocks_dbrepo.InsertUserStock"
	klogger.Enter(method)

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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return -1, err
	}

	klogger.Exit(method)
	return id, nil
}

// Function GetAllUserStocks returns all user stocks with names matching search if it is included and where t is after effective dt and before expiration date if it exists
func (m *PostgresDBRepo) GetAllUserStocks(userId int, stockType string, search string, t time.Time) ([]*models.UserStock, error) {
	method := "stocks_dbrepo.GetAllUserStocks"
	klogger.Enter(method)

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
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return usl, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
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
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

		recordCount = recordCount + 1
		usl = append(usl, &u)
	}

	klogger.Debug(method, "retrieved %d records", recordCount)
	klogger.Exit(method)
	return usl, nil
}

// Function GetAllUserStocksByDateRange returns all user stocks with names matching search if it is included and where the userStock was active during any part of the date range
func (m *PostgresDBRepo) GetAllUserStocksByDateRange(userId int, search string, ts time.Time, te time.Time) ([]*models.UserStock, error) {
	method := "stocks_dbrepo.GetAllUserStocks"
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
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return usl, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

	}

	defer rows.Close()

	for rows.Next() {
		var u models.UserStock
		err = rows.Scan(
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
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

		recordCount = recordCount + 1
		usl = append(usl, &u)
	}

	klogger.Debug(method, "retrieved %d records", recordCount)
	klogger.Exit(method)
	return usl, nil
}

// Function GetFirstUserStockBeforeDate returns the user stock with the closest effective date before or equal to d where userId and ticker match uId and t respectively
func (m *PostgresDBRepo) GetUserStockByUserIdTickerAndDate(uId int, t string, d time.Time) (models.UserStock, error) {
	method := "user_stocks_dbrepo.GetFirstUserStockBeforeDate"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			id, user_id, ticker, quantity, type, effective_dt, expiration_dt,
			create_dt, last_update_dt
		FROM user_stocks
		WHERE
			user_id = $1
		AND
			ticker = $2
		AND
			effective_dt <= $3
		AND
			(expiration_dt IS NULL OR expiration_dt >= $3)`

	row := m.DB.QueryRowContext(ctx, query, uId, t, d)

	var us models.UserStock

	err := row.Scan(
		&us.ID,
		&us.UserId,
		&us.Ticker,
		&us.Quantity,
		&us.Type,
		&us.EffectiveDt,
		&us.ExpirationDt,
		&us.CreateDt,
		&us.LastUpdateDt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return us, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return us, err
		}
	}

	klogger.Exit(method)
	return us, nil
}

// Function UpdateUserStock updates a user stock object in the database
func (m *PostgresDBRepo) UpdateUserStock(us models.UserStock) error {
	method := "stocks_dbrepo.UpdateStock"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`UPDATE user_stocks 
		SET
			quantity = $2,
			effective_dt = $3,
			expiration_dt = $4,
			type = $5,
			last_update_dt = $6
		WHERE
			id = $1`

	_, err := m.DB.ExecContext(ctx, stmt,
		us.ID,
		us.Quantity,
		us.EffectiveDt,
		us.ExpirationDt.Time,
		us.Type,
		time.Now(),
	)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	klogger.Exit(method)
	return nil
}
