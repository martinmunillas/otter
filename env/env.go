package env

import (
	"os"
	"strconv"

	"github.com/martinmunillas/otter/utils"
)

func RequiredStringEnvVar(key string) string {
	val := os.Getenv(key)
	if val == "" {
		utils.Throw("missing required env variable " + key)
	}
	return val
}

func OptionalStringEnvVar(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func OptionalBoolEnvVar(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val == "true"
}

func OptionalIntEnvVar(key string, defaultValue int64) int64 {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		utils.Throw("invalid int for env variable " + key)
	}
	return int64(v)
}
