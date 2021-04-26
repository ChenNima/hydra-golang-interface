package util

import "os"

// GetEnvWithDefault returns default value if key not exists in ENV
func GetEnvWithDefault(key string, defValue string) (res string) {
	if res = os.Getenv(key); res != "" {
		return res
	}
	return defValue
}
