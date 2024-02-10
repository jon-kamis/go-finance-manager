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

func (m *PostgresDBRepo) GetAllUserLoans(userId int, search string) ([]*models.Loan, error) {
	method := "loans_dbrepo.GetAllUserLoans"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query string
	var err error
	var rows *sql.Rows

	if search != "" {
		search = strings.ToLower(search)
		klogger.Debug(method, "searching for loans meeting criteria: %s", search)
		query = `
		SELECT
			id, user_id, loan_name, total_balance, total_cost, total_principal, total_interest, monthly_payment, interest_rate, loan_term,
			create_dt, last_update_dt
		FROM loans
		WHERE
			user_id = $1
			AND
			LOWER(loan_name) like '%' || $2 || '%'`
		rows, err = m.DB.QueryContext(ctx, query, userId, search)
	} else {
		query = `
		SELECT
			id, user_id, loan_name, total_balance, total_cost, total_principal, total_interest, monthly_payment, interest_rate, loan_term,
			create_dt, last_update_dt
		FROM loans
		WHERE
			user_id = $1`
		rows, err = m.DB.QueryContext(ctx, query, userId)
	}

	loans := []*models.Loan{}
	recordCount := 0

	if err != nil {
		if err == sql.ErrNoRows {
			return loans, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

	}

	defer rows.Close()

	for rows.Next() {
		var loan models.Loan
		err := rows.Scan(
			&loan.ID,
			&loan.UserID,
			&loan.Name,
			&loan.Total,
			&loan.TotalCost,
			&loan.TotalPayment,
			&loan.Interest,
			&loan.MonthlyPayment,
			&loan.InterestRate,
			&loan.LoanTerm,
			&loan.CreateDt,
			&loan.LastUpdateDt,
		)

		if err != nil {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return nil, err
		}

		recordCount = recordCount + 1
		loans = append(loans, &loan)
	}

	klogger.Debug(method, "retrieved %d records", recordCount)
	klogger.Exit(method)
	return loans, nil
}

func (m *PostgresDBRepo) GetLoanByID(id int) (models.Loan, error) {
	method := "loans_dbrepo.GetLoanByID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, user_id, loan_name, total_balance, total_cost, total_principal, total_interest, monthly_payment, interest_rate, loan_term,
		create_dt, last_update_dt
		FROM loans
		WHERE 
			id = $1`

	var loan models.Loan
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&loan.ID,
		&loan.UserID,
		&loan.Name,
		&loan.Total,
		&loan.TotalCost,
		&loan.TotalPayment,
		&loan.Interest,
		&loan.MonthlyPayment,
		&loan.InterestRate,
		&loan.LoanTerm,
		&loan.CreateDt,
		&loan.LastUpdateDt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			klogger.Info(method, constants.NoRowsReturnedMsg)
			klogger.Exit(method)
			return loan, nil
		} else {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return loan, err
		}

	}

	klogger.Exit(method)
	return loan, nil
}

func (m *PostgresDBRepo) DeleteLoanByID(id int) error {
	method := "loans_dbrepo.DeleteLoanByID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM loans
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

func (m *PostgresDBRepo) InsertLoan(loan models.Loan) (int, error) {
	method := "loans_dbrepo.InsertLoan"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`INSERT INTO loans 
			(user_id, loan_name, total_balance, total_cost, total_principal, total_interest, monthly_payment, interest_rate, loan_term, create_dt, last_update_dt)
		values 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt,
		loan.UserID,
		loan.Name,
		loan.Total,
		loan.TotalCost,
		loan.TotalPayment,
		loan.Interest,
		loan.MonthlyPayment,
		loan.InterestRate,
		loan.LoanTerm,
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

func (m *PostgresDBRepo) UpdateLoan(loan models.Loan) error {
	method := "loans_dbrepo.UpdateLoan"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt :=
		`UPDATE loans 
		SET
			loan_name = $2,
			total_balance = $3,
			total_cost = $4,
			total_principal = $5,
			total_interest = $6,
			monthly_payment = $7,
			interest_rate = $8,
			loan_term = $9,
			last_update_dt = $10
		WHERE
			id = $1`

	_, err := m.DB.ExecContext(ctx, stmt,
		loan.ID,
		loan.Name,
		loan.Total,
		loan.TotalCost,
		loan.TotalPayment,
		loan.Interest,
		loan.MonthlyPayment,
		loan.InterestRate,
		loan.LoanTerm,
		time.Now(),
	)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	klogger.Exit(method)
	return nil
}

func (m *PostgresDBRepo) DeleteLoansByUserID(id int) error {
	method := "loans_dbrepo.DeleteLoansByUserID"
	klogger.Enter(method)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE
		FROM loans
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
