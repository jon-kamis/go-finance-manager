package testingutils

import (
	"finance-manager-backend/cmd/finance-mngr/internal/authentication"
	"finance-manager-backend/cmd/finance-mngr/internal/config"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"time"
)

var TestAppConfig = config.GetDefaultConfig()

var TestAuth = authentication.Auth{
	Issuer:        config.GetEnvFromEnvValue(TestAppConfig.JWTIssuer),
	Audience:      config.GetEnvFromEnvValue(TestAppConfig.JWTAudience),
	Secret:        config.GetEnvFromEnvValue(TestAppConfig.JWTSecret),
	TokenExpiry:   time.Minute * 60,
	RefreshExpiry: time.Hour * 24,
	CookiePath:    "/",
	CookieName:    "Host-refresh_token",
	CookieDomain:  config.GetEnvFromEnvValue(TestAppConfig.CookieDomain),
}

func GetAdminJWT() (token string, err error) {
	method := "testing_auth.GetAdminJWT"
	fmlogger.Enter(method)

	u := authentication.JwtUser{
		ID:        TestingAdmin.ID,
		FirstName: TestingAdmin.FirstName,
		LastName:  TestingAdmin.LastName,
		Roles:     "admin",
	}

	//generate tokens
	tokens, err := TestAuth.GenerateTokenPair(&u)
	if err != nil {
		fmlogger.ExitError(method, "error occured when generating tokens", err)
		return "", err
	}

	fmlogger.Exit(method)
	return tokens.Token, nil

}

func GetUserJWT() (token string, err error) {
	method := "testing_auth.GetUserJWT"
	fmlogger.Enter(method)

	u := authentication.JwtUser{
		ID:        TestingUser.ID,
		FirstName: TestingUser.FirstName,
		LastName:  TestingUser.LastName,
		Roles:     "user",
	}

	//generate tokens
	tokens, err := TestAuth.GenerateTokenPair(&u)
	if err != nil {
		fmlogger.ExitError(method, "error occured when generating tokens", err)
		return "", err
	}

	fmlogger.Exit(method)
	return tokens.Token, nil

}
