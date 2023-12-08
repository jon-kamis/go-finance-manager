package main

import (
	"fmt"
	"os"
)

type env_value struct {
	envName    string
	defaultVal string
}

func getEnvFromEnvValue(env env_value) string {
	value := os.Getenv(env.envName)
	if value == "" {
		fmt.Printf("ussing fallback value for environment variable %s\n", env.envName)
		return env.defaultVal
	}
	fmt.Printf("loaded environment value for %s\n", env.envName)
	return value
}
