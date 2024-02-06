package dbrepo

import (
	test "finance-manager-backend/test"
	"finance-manager-backend/test/logtest"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jon-kamis/klogger"
)

var d PostgresDBRepo
var p test.DockerTestPlatform

func TestMain(m *testing.M) {
	logtest.SetKloggerTestFileNameEnv()

	method := "validation_test.TestMain"
	klogger.Enter(method)

	//Setup testing platform using docker
	p = test.Setup(m)
	d = PostgresDBRepo{DB: p.DB}

	//Execute Code
	code := m.Run()

	//Tear down testing platform
	test.TearDown(p)

	klogger.Exit(method)
	os.Exit(code)
}
