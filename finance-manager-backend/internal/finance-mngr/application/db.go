package application

import (
	"database/sql"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jon-kamis/klogger"
)

// Function OpenDB establishes a postgres database connection and returns a pointer to the database
func OpenDB(dsn string) (*sql.DB, error) {
	method := "db.OpenDB"
	klogger.Debug(method, "[ENTER]")
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	klogger.Debug(method, "[EXIT]")
	return db, nil
}

// Function ConnectToDB adds a Database connection to an Application object
func (app *Application) ConnectToDB() (*sql.DB, error) {
	method := "db.ConnectToDB"
	klogger.Enter(method)

	connection, err := OpenDB(app.DSN)
	if err != nil {
		klogger.ExitError(method, err.Error(), err)
		return nil, err
	}

	klogger.Info(method, "connected to postgres")
	klogger.Exit(method)
	return connection, nil
}
