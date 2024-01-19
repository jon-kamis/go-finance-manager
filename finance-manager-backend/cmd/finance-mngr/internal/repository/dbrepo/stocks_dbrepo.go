package dbrepo

import (
	"context"
	"database/sql"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"time"
)

func (m *PostgresDBRepo) InsertStock(s models.Stock) (int, error) {
	method := "stock_dbrepo.InsertStock"
	fmlogger.Enter(method)

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
		fmlogger.ExitError(method, "error occured when inserting new bill", err)
		return -1, err
	}

	fmlogger.Exit(method)
	return id, nil
}

func (m *PostgresDBRepo) InsertStockData(sl []models.Stock) error {
	method := "stock_dbrepo.InsertStockData"
	fmlogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`INSERT INTO stock_data 
			(ticker, high, low, open, close, date, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	var id int

	fmt.Printf("[%s] inserting %d records", method, len(sl))
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
			fmlogger.ExitError(method, "error occured when inserting stock data", err)
			return err
		}
	}

	fmlogger.Exit(method)
	return nil
}

func (m *PostgresDBRepo) GetStockByTicker(t string) (models.Stock, error) {
	method := "stocks_dbrepo.GetStockByTicker"
	fmlogger.Enter(method)

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
			fmlogger.Exit(method)
			return s, nil
		} else {
			fmlogger.ExitError(method, "database call returned with error", err)
			return s, err
		}
	}

	fmlogger.Exit(method)
	return s, nil
}

func (m *PostgresDBRepo) GetOldestStock() (models.Stock, error) {
	method := "stocks_dbrepo.GetOldestStock"
	fmlogger.Enter(method)

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
			fmlogger.Exit(method)
			return s, nil
		} else {
			fmlogger.ExitError(method, "database call returned with error", err)
			return s, err
		}
	}

	fmlogger.Exit(method)
	return s, nil
}

func (m *PostgresDBRepo) UpdateStock(s models.Stock) error {
	method := "stocks_dbrepo.UpdateStock"
	fmlogger.Enter(method)

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
		fmlogger.ExitError(method, "unexpected error occured when updating stock", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}
