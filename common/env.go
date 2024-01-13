package common

import (
	"os"
	"strconv"
)

func GetString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetBool(key string, defaultValue bool) bool {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		return defaultValue
	}

	value, errParseBool := strconv.ParseBool(rawValue)
	if errParseBool != nil {
		return defaultValue
	}

	return value
}

func GetInt32(key string, defaultValue int32) int32 {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		return defaultValue
	}

	value, errParseInt := strconv.ParseInt(rawValue, 10, 32)
	if errParseInt != nil {
		return defaultValue
	}

	return int32(value)
}
