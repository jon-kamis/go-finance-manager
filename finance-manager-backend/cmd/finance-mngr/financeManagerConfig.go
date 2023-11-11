package main

type FinanceManagerConfig struct {
	DSN          env_value
	JWTSecret    env_value
	JWTIssuer    env_value
	JWTAudience  env_value
	CookieDomain env_value
	Domain       env_value
	FrontendUrl  env_value
}

type env_value struct {
	envName    string
	defaultVal string
}
