package utils

import (
	"fmt"
	"os"
	"strings"
)

const EnvRealHostKey = "APP_REAL_HOST_%s"

func GetRealHostFromEnv(key string) string {
	envKey := fmt.Sprintf(EnvRealHostKey, strings.ToUpper(key))
	return os.Getenv(envKey)
}
