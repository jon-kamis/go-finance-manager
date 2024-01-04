package handlers

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/repository/dbrepo"
	"finance-manager-backend/cmd/finance-mngr/internal/testingutils"
	"finance-manager-backend/cmd/finance-mngr/internal/validation"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var fmh FinanceManagerHandler

func TestMain(m *testing.M) {
	method := "validation_test.TestMain"
	fmlogger.Enter(method)

	//Setup testing platform using docker
	p := testingutils.Setup(m)
	db := &dbrepo.PostgresDBRepo{DB: p.DB}
	fmh = FinanceManagerHandler{
		DB:        db,
		Validator: &validation.FinanceManagerValidator{DB: db},
		Auth: testingutils.TestAuth,
	}

	//Execute Code
	code := m.Run()

	//Tear down testing platform
	testingutils.TearDown(p)

	fmlogger.Exit(method)
	os.Exit(code)
}
