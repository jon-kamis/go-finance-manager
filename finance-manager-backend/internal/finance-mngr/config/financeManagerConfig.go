//Package config contains files used to read ENV values into the application as well as providing default values for each env variable
package config

//Type FinanceManagerConfig contains all of the ENV values used by the application
type FinanceManagerConfig struct {
	DSN          Env_value
	JWTSecret    Env_value
	JWTIssuer    Env_value
	JWTAudience  Env_value
	CookieDomain Env_value
	Domain       Env_value
	FrontendUrl  Env_value
	TimeZone     Env_value
}

//Function GetDefaultConfig returns a FinanceManagerConfig object containing the default values for each environment variable
func GetDefaultConfig() FinanceManagerConfig {
	config := FinanceManagerConfig{
		DSN: Env_value{
			envName:    "DSN",
			defaultVal: "host=localhost port=5432 user=postgres password=postgres dbname=financemanager sslmode=disable timezone=UTC connect_timeout=5",
		},
		JWTSecret: Env_value{
			envName:    "JWTSecret",
			defaultVal: "verysecret",
		},
		JWTIssuer: Env_value{
			envName:    "JWTIssuer",
			defaultVal: "fm.com",
		},
		JWTAudience: Env_value{
			envName:    "JWTAudience",
			defaultVal: "fm.com",
		},
		CookieDomain: Env_value{
			envName:    "CookieDomain",
			defaultVal: "localhost",
		},
		Domain: Env_value{
			envName:    "Domain",
			defaultVal: "fm.com",
		},
		FrontendUrl: Env_value{
			envName:    "FrontendUrl",
			defaultVal: "http://localhost:3000",
		},
		TimeZone: Env_value{
			envName:    "TimeZone",
			defaultVal: "America/New_York",
		},
	}

	return config
}
