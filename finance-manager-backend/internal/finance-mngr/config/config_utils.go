package config

import (
	"finance-manager-backend/internal/finance-mngr/constants"
	"fmt"
	"os"

	"github.com/jon-kamis/klogger"
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
	klogger.Trace(method, constants.EnterLog)
	value := os.Getenv(env.envName)
	
	if value == "" {
		pwd, _ := os.Getwd()
		propFile, err := os.ReadFile(pwd + constants.PropertyFileName)

		if err == nil {
			var props propData
			err = yaml.Unmarshal(propFile, &props)

			if err != nil {
				klogger.Error(method, "error Reading property file")
			} else if props.Config[env.envName] != "" {
				klogger.Trace(method, "loaded property file value for property: %s", env.envName)
				return fmt.Sprintf("%v", props.Config[env.envName])
			}
		}

		klogger.Trace(method, "used default value for property: %s", env.envName)
		klogger.Trace(method, constants.ExitLog)
		return env.defaultVal
	}
	klogger.Trace(method, "loaded env value for property: %s", env.envName)
	klogger.Trace(method, constants.ExitLog)
	return value
}
