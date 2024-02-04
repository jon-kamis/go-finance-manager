package application

import (
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/fmlogger"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Function OpenDB establishes a postgres database connection and returns a pointer to the database
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

// Function ConnectToDB adds a Database connection to an Application object
func (app *Application) ConnectToDB() (*sql.DB, error) {
	method := "db.ConnectToDB"
	fmlogger.Enter(method)

	connection, err := OpenDB(app.DSN)
	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		return nil, err
	}

	fmlogger.Info(method, "connected to postgres")
	fmlogger.Exit(method)
	return connection, nil
}
