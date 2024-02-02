package fmservice

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/repository/dbrepo"
	"finance-manager-backend/test"
	"os"
	"testing"
)

var p test.DockerTestPlatform
var fms FMService

func TestMain(m *testing.M) {
	method := "validation_test.TestMain"
	fmlogger.Enter(method)

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

	fmlogger.Exit(method)
	os.Exit(code)
}
