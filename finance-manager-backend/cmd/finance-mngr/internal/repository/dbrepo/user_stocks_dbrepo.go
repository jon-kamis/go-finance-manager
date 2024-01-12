package dbrepo

import (
	"context"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
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
