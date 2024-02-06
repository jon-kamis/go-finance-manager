package fmservice

import (
	"finance-manager-backend/internal/finance-mngr/repository/dbrepo"
	"finance-manager-backend/test"
	"finance-manager-backend/test/logtest"
	"os"
	"testing"

	"github.com/jon-kamis/klogger"
)

var p test.DockerTestPlatform
var fms FMService

func TestMain(m *testing.M) {
	logtest.SetKloggerTestFileNameEnv()

	method := "validation_test.TestMain"
	klogger.Enter(method)

	//Setup testing platform using docker
	p = test.Setup(m)
	db := &dbrepo.PostgresDBRepo{DB: p.DB}
	fms = FMService{
		DB: db,
	}

	//Execute Code
	code := m.Run()

	//Tear down testing platform
	test.TearDown(p)

	klogger.Exit(method)
	os.Exit(code)
}
