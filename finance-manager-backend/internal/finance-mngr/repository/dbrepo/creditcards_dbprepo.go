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

func (m *PostgresDBRepo) GetAllUserCreditCards(userId int, search string) ([]*models.CreditCard, error) {
	method := "creditcards.GetAllUserCreditCards"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query string
	var err error
	var rows *sql.Rows

	if search != "" {
		search = strings.ToLower(search)
		klogger.Debug("Searching for creditcards meeting criteria: %s", search)
		query = `
		SELECT
			id, user_id, name, balance, credit_limit, apr, min_pay, min_pay_percentage,
			create_dt, last_update_dt
		FROM credit_cards
		WHERE
			user_id = $1
			AND
			LOWER(name) like '%' || $2 || '%'`
		rows, err = m.DB.QueryContext(ctx, query, userId, search)
	} else {
		query = `
		SELECT
		id, user_id, name, balance, credit_limit, apr, min_pay, min_pay_percentage,
			create_dt, last_update_dt
		FROM credit_cards
		WHERE
			user_id = $1`
		rows, err = m.DB.QueryContext(ctx, query, userId)
	}

	creditcards := []*models.CreditCard{}
	recordCount := 0

	if err != nil {
		if err == sql.ErrNoRows {
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return creditcards, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

	}

	defer rows.Close()

	for rows.Next() {
		var cc models.CreditCard
		err := rows.Scan(
			&cc.ID,
			&cc.UserID,
			&cc.Name,
			&cc.Balance,
			&cc.Limit,
			&cc.APR,
			&cc.MinPayment,
			&cc.MinPaymentPercentage,
			&cc.CreateDt,
			&cc.LastUpdateDt,
		)

		if err != nil {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

		recordCount = recordCount + 1
		creditcards = append(creditcards, &cc)
	}

	klogger.Debug(method, "retrieved %d records\n", recordCount)
	klogger.Exit(method)
	return creditcards, nil
}

func (m *PostgresDBRepo) GetCreditCardByID(id int) (models.CreditCard, error) {
	method := "creditcards_dbrepo.GetCreditCardsByID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, user_id, name, balance, credit_limit, apr, min_pay, min_pay_percentage,
			create_dt, last_update_dt
		FROM credit_cards
		WHERE 
			id = $1`

	var cc models.CreditCard
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&cc.ID,
		&cc.UserID,
		&cc.Name,
		&cc.Balance,
		&cc.Limit,
		&cc.APR,
		&cc.MinPayment,
		&cc.MinPaymentPercentage,
		&cc.CreateDt,
		&cc.LastUpdateDt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return cc, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return cc, err
		}
	}

	klogger.Exit(method)
	return cc, nil
}

func (m *PostgresDBRepo) InsertCreditCard(cc models.CreditCard) (int, error) {
	method := "creditcards_dbrepo.InsertCreditCard"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`INSERT INTO credit_cards 
			(user_id, name, balance, credit_limit, apr, min_pay, min_pay_percentage, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt,
		cc.UserID,
		cc.Name,
		cc.Balance,
		cc.Limit,
		cc.APR,
		cc.MinPayment,
		cc.MinPaymentPercentage,
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

func (m *PostgresDBRepo) DeleteCreditCardsByUserID(id int) error {
	method := "creditcardss_dbrepo.DeleteCreditCardsByUserID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM credit_cards
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

func (m *PostgresDBRepo) DeleteCreditCardsByID(id int) error {
	method := "creditcardss_dbrepo.DeleteCreditCardsByUserID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM credit_cards
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

func (m *PostgresDBRepo) UpdateCreditCard(cc models.CreditCard) error {
	method := "creditcards_dbrepo.UpdateCreditCard"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`UPDATE credit_cards 
		SET
			user_id = $2,
			name = $3,
			balance = $4,
			credit_limit = $5,
			apr = $6,
			min_pay = $7,
			min_pay_percentage = $8,
			create_dt = $9,
			last_update_dt = $10
		WHERE
			id = $1`

	_, err := m.DB.ExecContext(ctx, stmt,
		cc.ID,
		cc.UserID,
		cc.Name,
		cc.Balance,
		cc.Limit,
		cc.APR,
		cc.MinPayment,
		cc.MinPaymentPercentage,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	klogger.Exit(method)
	return nil
}
