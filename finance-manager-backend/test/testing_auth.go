package test

import (
	"finance-manager-backend/internal/finance-mngr/authentication"
	"finance-manager-backend/internal/finance-mngr/config"
	"testing"
	"time"

	"github.com/jon-kamis/klogger"
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
	klogger.Enter(method)

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

	klogger.Exit(method)
	return tokens.Token

}

func GetUserJWT(t *testing.T) (token string) {
	method := "testing_auth.GetUserJWT"
	klogger.Enter(method)

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

	klogger.Exit(method)
	return tokens.Token

}

func GetUserJWTWithId(t *testing.T, id int) (token string) {
	method := "testing_auth.GetUserJWT"
	klogger.Enter(method)

	u := authentication.JwtUser{
		ID:        id,
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

	klogger.Exit(method)
	return tokens.Token

}
