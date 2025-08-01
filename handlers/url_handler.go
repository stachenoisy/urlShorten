package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"urlShort/config"
	"urlShort/storage"
	"urlShort/utils"
)

// Global değişkenler
var AppConfig *config.Config
var AppStorage storage.Storage

// Config'i set eder
func SetConfig(cfg *config.Config) {
	AppConfig = cfg
}

// Storage'ı set eder
func SetStorage(s storage.Storage) {
	AppStorage = s
}

// Ana sayfa
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "pages/index.html")
		return
	}
}

// URL kısaltma işlemi
func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Sadece POST metoduna izin verilir", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		URL string `json:"url"`
	}

	// JSON verisi parse et
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
		return
	}

	// URL validation
	if !utils.IsValidURL(requestData.URL) {
		http.Error(w, "Geçersiz URL formatı", http.StatusBadRequest)
		return
	}

	// Kısa kod üret
	shortCode := utils.GenerateShortCode(AppConfig.ShortenerLength)

	// URL objesini oluştur
	url := &storage.URL{
		Original:  requestData.URL,
		Short:     shortCode,
		CreatedAt: time.Now(),
		Clicks:    0,
	}

	// Storage'a kaydet
	if err := AppStorage.Save(url); err != nil {
		http.Error(w, "URL kaydedilemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Yanıt gönder
	response := map[string]string{
		"original":   url.Original,
		"short_url":  fmt.Sprintf("%s/s/%s", AppConfig.BaseURL, url.Short),
		"short_code": url.Short,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Kısa URL'den orijinal URL'e yönlendirme
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// URL'den kısa kodu al (/s/abc123 -> abc123)
	shortCode := strings.TrimPrefix(r.URL.Path, "/s/")

	// Storage'dan URL'i bul
	url, err := AppStorage.Get(shortCode)
	if err != nil {
		http.Error(w, "URL bulunamadı", http.StatusNotFound)
		return
	}

	// Tıklanma sayısını artır
	if err := AppStorage.IncrementClicks(shortCode); err != nil {
		fmt.Printf("Tıklanma sayısı güncellenemedi: %v\n", err)
	}

	// Orijinal URL'e yönlendir
	http.Redirect(w, r, url.Original, http.StatusMovedPermanently)
}

// URL istatistikleri
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/stats/")

	// Storage'dan URL'yi bul
	url, err := AppStorage.Get(shortCode)
	if err != nil {
		http.Error(w, "URL bulunamadı", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(url)
}

// Tüm URL'leri listeler
func ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Sadece GET metoduna izin verilir", http.StatusMethodNotAllowed)
		return
	}

	// Tüm URL'leri getir
	urls, err := AppStorage.GetAll()
	if err != nil {
		http.Error(w, "URL'ler getirilemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Yanıt gönder
	response := map[string]interface{}{
		"count": len(urls),
		"urls":  urls,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
