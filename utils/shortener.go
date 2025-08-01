package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

// Kısa kod üretir
func GenerateShortCode(ShortenerLength int) string {
	bytes := make([]byte, ShortenerLength)
	rand.Read(bytes)
	encoded := base64.URLEncoding.EncodeToString(bytes)
	return strings.ReplaceAll(strings.ReplaceAll(encoded[:ShortenerLength], "-", ""), "_", "")
}

// Basit bir şekilde URL doğrulaması yapar
func IsValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
