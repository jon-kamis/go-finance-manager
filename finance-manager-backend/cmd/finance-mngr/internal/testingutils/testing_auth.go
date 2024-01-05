package testingutils

import (
	"finance-manager-backend/cmd/finance-mngr/internal/authentication"
	"finance-manager-backend/cmd/finance-mngr/internal/config"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"testing"
	"time"
)

func GetTestAuth() authentication.Auth {
	return authentication.Auth{
		Issuer:        config.GetEnvFromEnvValue(TestAppConfig.JWTIssuer),
		Audience:      config.GetEnvFromEnvValue(TestAppConfig.JWTAudience),
		Secret:        config.GetEnvFromEnvValue(TestAppConfig.JWTSecret),
		TokenExpiry:   time.Minute * 60,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "Host-refresh_token",
		CookieDomain:  config.GetEnvFromEnvValue(TestAppConfig.CookieDomain),
	}
}

func GetAdminJWT(t *testing.T) (token string) {
	method := "testing_auth.GetAdminJWT"
	fmlogger.Enter(method)

	u := authentication.JwtUser{
		ID:        TestingAdmin.ID,
		FirstName: TestingAdmin.FirstName,
		LastName:  TestingAdmin.LastName,
		Roles:     "admin",
	}

	//generate tokens
	auth := GetTestAuth()
	tokens, err := auth.GenerateTokenPair(&u)
	if err != nil {
		t.Errorf("error occured when generating tokens: %s", err)
	}

	fmlogger.Exit(method)
	return tokens.Token

}

func GetUserJWT(t *testing.T) (token string) {
	method := "testing_auth.GetUserJWT"
	fmlogger.Enter(method)

	u := authentication.JwtUser{
		ID:        TestingUser.ID,
		FirstName: TestingUser.FirstName,
		LastName:  TestingUser.LastName,
		Roles:     "user",
	}

	//generate tokens
	auth := GetTestAuth()
	tokens, err := auth.GenerateTokenPair(&u)

	if err != nil {
		t.Errorf("error occured when generating tokens: %s", err)
	}

	fmlogger.Exit(method)
	return tokens.Token

}
