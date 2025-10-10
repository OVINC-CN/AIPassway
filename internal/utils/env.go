package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/OVINC-CN/AIPassway/internal/logger"
)

const EnvRealHostKey = "APP_REAL_HOST_%s"

func GetRealHostFromEnv(key string) string {
	envKey := fmt.Sprintf(EnvRealHostKey, strings.ToUpper(key))
	return os.Getenv(envKey)
}

func GetConfigIntFromEnv(key string, defaultValue int) int {
	envValue := os.Getenv(key)
	if envValue == "" {
		return defaultValue
	}

	var intValue int
	_, err := fmt.Sscanf(envValue, "%d", &intValue)
	if err != nil {
		logger.Logger().Warnf("[GetConfigIntFromEnv] parse int from env failed: %s=%s\nerror: %v", key, envValue, err)
		return defaultValue
	}

	return intValue
}
