package restmodels

import (
	"finance-manager-backend/test/logtest"
	"os"
	"testing"

	"github.com/jon-kamis/klogger"
)

func TestMain(m *testing.M) {
	logtest.SetKloggerTestFileNameEnv()

	method := "restmodels_test.TestMain"
	klogger.Enter(method)

	//Execute Code
	code := m.Run()

	klogger.Exit(method)
	os.Exit(code)
}
