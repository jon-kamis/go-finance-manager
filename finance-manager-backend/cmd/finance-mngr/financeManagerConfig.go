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

func getDefaultConfig() FinanceManagerConfig {
	config := FinanceManagerConfig{
		DSN: env_value{
			envName:    "DSN",
			defaultVal: "host=localhost port=5432 user=postgres password=postgres dbname=financemanager sslmode=disable timezone=UTC connect_timeout=5",
		},
		JWTSecret: env_value{
			envName:    "JWTSecret",
			defaultVal: "verysecret",
		},
		JWTIssuer: env_value{
			envName:    "JWTIssuer",
			defaultVal: "fm.com",
		},
		JWTAudience: env_value{
			envName:    "JWTAudience",
			defaultVal: "fm.com",
		},
		CookieDomain: env_value{
			envName:    "CookieDomain",
			defaultVal: "localhost",
		},
		Domain: env_value{
			envName:    "Domain",
			defaultVal: "fm.com",
		},
		FrontendUrl: env_value{
			envName:    "FrontendUrl",
			defaultVal: "http://localhost:3000",
		},
	}

	return config
}
