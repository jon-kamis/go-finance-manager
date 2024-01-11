package config

import (
	"fmt"
	"os"
)

//Type Env_value is a pairing of environment variables and their default values
type Env_value struct {
	envName    string
	defaultVal string
}

//Function GetEnvFromEnvValue will attempt to read in the envrionment variable from name
//if one was not provided it instead returns the default value associated with an Env_value object
func GetEnvFromEnvValue(env Env_value) string {
	value := os.Getenv(env.envName)
	if value == "" {
		fmt.Printf("using fallback value for environment variable %s\n", env.envName)
		return env.defaultVal
	}
	fmt.Printf("loaded environment value for %s\n", env.envName)
	return value
}
