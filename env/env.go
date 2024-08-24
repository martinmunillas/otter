package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/martinmunillas/otter/utils"
)

func RequiredStringEnvVar(key string) string {
	val := os.Getenv(key)
	if val == "" {
		utils.Throw(fmt.Sprintf("missing required env variable `%s`", key))
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
		utils.Throw(fmt.Sprintf("invalid int \"%s\" for env variable `%s`", val, key))
	}
	return int64(v)
}

func RequiredIntEnvVar(key string) int64 {
	val := os.Getenv(key)
	if val == "" {
		utils.Throw(fmt.Sprintf("missing required int for env variable `%s`", key))
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		utils.Throw(fmt.Sprintf("invalid required int \"%s\" for env variable `%s`", val, key))
	}
	return int64(v)
}
