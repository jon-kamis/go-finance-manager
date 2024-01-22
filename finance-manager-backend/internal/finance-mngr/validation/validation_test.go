package validation

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/repository/dbrepo"
	"finance-manager-backend/test"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var fmv FinanceManagerValidator
var p test.DockerTestPlatform

func TestMain(m *testing.M) {
	method := "validation_test.TestMain"
	fmlogger.Enter(method)

	//Setup testing platform using docker
	p = test.Setup(m)
	fmv = FinanceManagerValidator{DB: &dbrepo.PostgresDBRepo{DB: p.DB}}

	//Execute Code
	code := m.Run()

	//Tear down testing platform
	test.TearDown(p)

	fmlogger.Exit(method)
	os.Exit(code)
}
