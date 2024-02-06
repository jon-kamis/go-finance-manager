package authentication

import (
	"finance-manager-backend/internal/finance-mngr/config"
	"finance-manager-backend/test/logtest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jon-kamis/klogger"
)

var testAppConfig = config.GetDefaultConfig()
var auth = Auth{
	Issuer:        config.GetEnvFromEnvValue(testAppConfig.JWTIssuer),
	Audience:      config.GetEnvFromEnvValue(testAppConfig.JWTAudience),
	Secret:        config.GetEnvFromEnvValue(testAppConfig.JWTSecret),
	TokenExpiry:   time.Minute * 60,
	RefreshExpiry: time.Hour * 24,
	CookiePath:    "/",
	CookieName:    "Host-refresh_token",
	CookieDomain:  config.GetEnvFromEnvValue(testAppConfig.CookieDomain),
}

var usr = JwtUser{
	ID:        1,
	FirstName: "usr",
	LastName:  "usr",
	Roles:     "admin",
}

func TestMain(m *testing.M) {
	logtest.SetKloggerTestFileNameEnv()

	method := "auth_test.TestMain"
	klogger.Enter(method)

	//Execute Code
	code := m.Run()

	klogger.Exit(method)
	os.Exit(code)
}

func TestGenerateTokenPairAndParseAndVerifyToken(t *testing.T) {
	method := "auth_test.TestGenerateTokenPair"
	klogger.Enter(method)

	tokenPairs, err := auth.GenerateTokenPair(&usr)

	if err != nil {
		t.Errorf("unexpected error generating tokens %s", err)
	}

	if tokenPairs.Token == "" || tokenPairs.RefreshToken == "" {
		t.Errorf("token values returned empty %s", err)
	}

	_, claims, err := auth.ParseAndVerifyToken(tokenPairs.Token)

	if err != nil {
		t.Errorf("unexpected error parsing tokens %s", err)
	}

	id, err := strconv.Atoi(claims.Subject)

	if err != nil {
		t.Errorf("unexpected error decoding claims %s", err)
	}

	if id != usr.ID {
		t.Errorf("decoded claims subject does not match original value")
	}

	klogger.Exit(method)
}
