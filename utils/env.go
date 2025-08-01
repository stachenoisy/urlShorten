package utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// LoadEnvFile - .env dosyasını yükler
func LoadEnvFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err // .env dosyası yoksa hata döndür
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Boş satırları ve yorumları atla
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// KEY=VALUE formatını parse et
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Tırnak işaretlerini temizle
			value = strings.Trim(value, `"'`)

			// Environment variable'ı set et
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// Environment variable'ı alır, yoksa default değer döner (Type string)
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Environment variable'ı alır, yoksa default değer döner (Type int)
func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Environment variable'ı alır, yoksa default değer döner (Type bool)
func GetEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
