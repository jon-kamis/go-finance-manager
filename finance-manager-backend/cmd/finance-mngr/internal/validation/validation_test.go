package validation

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/repository/dbrepo"
	"finance-manager-backend/cmd/finance-mngr/internal/testingutils"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var fmv FinanceManagerValidator

func TestMain(m *testing.M) {
	method := "validation_test.TestMain"
	fmlogger.Enter(method)

	//Setup testing platform using docker
	p := testingutils.Setup(m)
	fmv = FinanceManagerValidator{DB: &dbrepo.PostgresDBRepo{DB: p.DB}}

	//Execute Code
	code := m.Run()

	//Tear down testing platform
	testingutils.TearDown(p)

	fmlogger.Exit(method)
	os.Exit(code)
}
