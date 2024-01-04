package config

import (
	"fmt"
	"os"
)

type Env_value struct {
	envName    string
	defaultVal string
}

func GetEnvFromEnvValue(env Env_value) string {
	value := os.Getenv(env.envName)
	if value == "" {
		fmt.Printf("ussing fallback value for environment variable %s\n", env.envName)
		return env.defaultVal
	}
	fmt.Printf("loaded environment value for %s\n", env.envName)
	return value
}
