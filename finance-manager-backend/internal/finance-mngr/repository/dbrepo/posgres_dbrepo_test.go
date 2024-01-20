package dbrepo

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/testingutils"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var d PostgresDBRepo
var p testingutils.DockerTestPlatform

func TestMain(m *testing.M) {
	method := "validation_test.TestMain"
	fmlogger.Enter(method)

	//Setup testing platform using docker
	p = testingutils.Setup(m)
	d = PostgresDBRepo{DB: p.DB}

	//Execute Code
	code := m.Run()

	//Tear down testing platform
	testingutils.TearDown(p)

	fmlogger.Exit(method)
	os.Exit(code)
}
