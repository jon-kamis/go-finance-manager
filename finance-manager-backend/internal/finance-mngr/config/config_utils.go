package config

import (
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Type Env_value is a pairing of environment variables and their default values
type Env_value struct {
	envName    string
	defaultVal string
}

type propData struct {
	Config map[string]string
}

// Function GetEnvFromEnvValue will attempt to read in the envrionment variable from name
// if one was not provided it then checks a property file and if that is also not found
// returns the default value associated with an Env_value object
func GetEnvFromEnvValue(env Env_value) string {
	method := "config_utils.GetEnvFromEnvValue"
	fmlogger.Enter(method)

	value := os.Getenv(env.envName)
	if value == "" {
		pwd, _ := os.Getwd()
		propFile, err := os.ReadFile(pwd + constants.PropertyFileName)

		if err != nil {
			fmlogger.Info(method, "property file not found")
		} else {
			var props propData
			err = yaml.Unmarshal(propFile, &props)

			if err != nil {
				fmlogger.Info(method, "error Reading property file")
			} else if props.Config[env.envName] != "" {
				fmlogger.Info(method, "loaded value from properties file for %s", env.envName)
				fmlogger.Exit(method)
				return fmt.Sprintf("%v",props.Config[env.envName])
			}
		}

		fmlogger.Info(method, "using fallback value for environment variable %s", env.envName)
		fmlogger.Exit(method)
		return env.defaultVal
	}
	
	fmlogger.Info(method, "loaded environment value for %s\n", env.envName)
	fmlogger.Exit(method)
	return value
}
