package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"fmt"
	"time"

	"github.com/jon-kamis/klogger"
)

func (m *PostgresDBRepo) InsertStock(s models.Stock) (int, error) {
	method := "stock_dbrepo.InsertStock"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`INSERT INTO stocks 
			(ticker, high, low, open, close, date, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt,
		s.Ticker,
		s.High,
		s.Low,
		s.Open,
		s.Close,
		s.Date,
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

func (m *PostgresDBRepo) InsertStockData(sl []models.Stock) error {
	method := "stock_dbrepo.InsertStockData"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	foundErr := false
	var errList []error

	stmt :=
		`INSERT INTO stock_data 
			(ticker, high, low, open, close, date, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	var id int

	klogger.Debug(method, "inserting %d records", len(sl))
	for _, s := range sl {
		err := m.DB.QueryRowContext(ctx, stmt,
			s.Ticker,
			s.High,
			s.Low,
			s.Open,
			s.Close,
			s.Date,
			time.Now(),
			time.Now(),
		).Scan(&id)

		if err != nil {

			foundErr = true
			errMsg := fmt.Sprintf("Error occured when inserting %s for %v: %s", s.Ticker, s.Date, err.Error())
			errList = append(errList, errors.New(errMsg))
		}
	}

	if foundErr {
		for _, err := range errList {
			klogger.Error(method, err.Error())
		}

		err := errors.New(constants.InsertMultStockDataError)
		klogger.ExitError(method, err.Error())
		return err
	}

	klogger.Exit(method)
	return nil
}

func (m *PostgresDBRepo) GetLatestStockDataByTicker(t string) (models.Stock, error) {
	method := "stocks_dbrepo.GetLatestStockDataByTicker"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, ticker, high, low, open, close, date,
			create_dt, last_update_dt
		FROM stock_data
		WHERE ticker = $1
		ORDER BY  date desc
		limit 1`

	var s models.Stock
	row := m.DB.QueryRowContext(ctx, query, t)

	err := row.Scan(
		&s.ID,
		&s.Ticker,
		&s.High,
		&s.Low,
		&s.Open,
		&s.Close,
		&s.Date,
		&s.CreateDt,
		&s.LastUpdateDt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return s, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return s, err
		}
	}

	klogger.Exit(method)
	return s, nil
}

func (m *PostgresDBRepo) GetStockDataByTickerAndDateRange(t string, sd time.Time, ed time.Time) ([]models.Stock, error) {
	method := "stocks_dbrepo.GetStockDataByTickerAndDateRange"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, ticker, high, low, open, close, date,
			create_dt, last_update_dt
		FROM stock_data
		WHERE ticker = $1
		AND date >= $2
		AND date <= $3
		ORDER BY  date asc`

	rows, err := m.DB.QueryContext(ctx, query, t, sd, ed)

	stocks := []models.Stock{}
	recordCount := 0

	if err != nil {
		if err == sql.ErrNoRows {
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return stocks, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

	}

	defer rows.Close()

	for rows.Next() {
		var s models.Stock
		err := rows.Scan(
			&s.ID,
			&s.Ticker,
			&s.High,
			&s.Low,
			&s.Open,
			&s.Close,
			&s.Date,
			&s.CreateDt,
			&s.LastUpdateDt,
		)

		if err != nil {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

		recordCount = recordCount + 1
		stocks = append(stocks, s)
	}

	klogger.Exit(method)
	return stocks, nil
}

func (m *PostgresDBRepo) GetStockByTicker(t string) (models.Stock, error) {
	method := "stocks_dbrepo.GetStockByTicker"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, ticker, high, low, open, close, date,
			create_dt, last_update_dt
		FROM stocks
		WHERE 
			ticker = $1`

	var s models.Stock
	row := m.DB.QueryRowContext(ctx, query, t)

	err := row.Scan(
		&s.ID,
		&s.Ticker,
		&s.High,
		&s.Low,
		&s.Open,
		&s.Close,
		&s.Date,
		&s.CreateDt,
		&s.LastUpdateDt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return s, nil
		} else {
			klogger.ExitError(method, "database call returned with error", err)
			return s, err
		}
	}

	klogger.Exit(method)
	return s, nil
}

func (m *PostgresDBRepo) GetOldestStock() (models.Stock, error) {
	method := "stocks_dbrepo.GetOldestStock"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, ticker, high, low, open, close, date,
			create_dt, last_update_dt
		FROM stocks
		ORDER BY  date, last_update_dt asc
		limit 1`

	var s models.Stock
	row := m.DB.QueryRowContext(ctx, query)

	err := row.Scan(
		&s.ID,
		&s.Ticker,
		&s.High,
		&s.Low,
		&s.Open,
		&s.Close,
		&s.Date,
		&s.CreateDt,
		&s.LastUpdateDt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return s, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return s, err
		}
	}

	klogger.Exit(method)
	return s, nil
}

func (m *PostgresDBRepo) UpdateStock(s models.Stock) error {
	method := "stocks_dbrepo.UpdateStock"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`UPDATE stocks 
		SET
			high = $2,
			low = $3,
			open = $4,
			close = $5,
			date = $6,
			last_update_dt = $7
		WHERE
			id = $1`

	_, err := m.DB.ExecContext(ctx, stmt,
		s.ID,
		s.High,
		s.Low,
		s.Open,
		s.Close,
		s.Date,
		time.Now(),
	)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	klogger.Exit(method)
	return nil
}
