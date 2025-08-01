package config

import (
	"fmt"
	"urlShort/utils"
)

// Config - Type tanımı
type Config struct {
	Port            string
	Host            string
	BaseURL         string
	SSLenabled      bool
	Database        string
	Debug           bool
	ShortenerLength int
}

// Environment üzerinden config yükler
func Load() *Config {
	// .env dosyasını yükle (hata varsa sessizce devam et)
	if err := utils.LoadEnvFile(".env"); err != nil {
		// .env dosyası yoksa veya okunamazsa, sadece environment variable'ları kullan
		fmt.Printf("⚠️  .env dosyası yüklenemedi (opsiyonel): %v\n", err)
	}

	port := utils.GetEnv("PORT", "8080")
	host := utils.GetEnv("HOST", "localhost")
	database := utils.GetEnv("DATABASE", "memory")
	debug := utils.GetEnvAsBool("DEBUG", false)
	shortenerlength := utils.GetEnvAsInt("SHORTENER_LENGTH", 6)
	sslenabled := utils.GetEnvAsBool("SSL", false)

	var baseURL string
	if sslenabled {
		baseURL = fmt.Sprintf("https://%s:%s", host, port)
	} else {
		baseURL = fmt.Sprintf("http://%s:%s", host, port)
	}

	return &Config{
		Port:            port,
		Host:            host,
		BaseURL:         baseURL,
		Database:        database,
		Debug:           debug,
		ShortenerLength: shortenerlength,
		SSLenabled:      sslenabled,
	}
}

// Server address döner
func (c *Config) GetAddress() string {
	return ":" + c.Port
}

// Tam URL döner
func (c *Config) GetFullURL(path string) string {
	return c.BaseURL + path
}
