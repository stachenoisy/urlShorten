package main

import (
	"fmt"
	"log"
	"net/http"
	"urlShort/config"
	"urlShort/handlers"
	"urlShort/storage"
)

func main() {
	// Config yükle
	cfg := config.Load()
	handlers.SetConfig(cfg)

	// Storage'ı başlat
	appStorage, err := storage.NewStorage(cfg.Database)
	if err != nil {
		log.Fatalf("Storage başlatılamadı: %v", err)
	}
	defer appStorage.Close()

	handlers.SetStorage(appStorage)

	// Route'ları tanımla
	http.HandleFunc("/", handlers.HomeHandler)           // Ana sayfa
	http.HandleFunc("/shorten", handlers.ShortenHandler) // URL kısaltma
	http.HandleFunc("/s/", handlers.RedirectHandler)     // Yönlendirme
	http.HandleFunc("/stats/", handlers.StatsHandler)    // İstatistikler
	if cfg.Debug {
		http.HandleFunc("/list", handlers.ListHandler) // Tüm URL'leri listele
	}

	// Sistem bilgileri
	fmt.Println("🚀 URL Kısaltma Servisi çalışıyor...")
	fmt.Printf("📍 Ana sayfa: %s\n", cfg.BaseURL)
	fmt.Printf("⚙️  Port: %s\n", cfg.Port)
	fmt.Printf("🖥️  Host: %s\n", cfg.Host)
	fmt.Printf("🗄️  Database: %s\n", cfg.Database)
	if cfg.Debug {
		fmt.Println("🐛 Debug modu: AÇIK")
		fmt.Printf("📊 Liste endpoint: %s/list\n", cfg.BaseURL)
	}

	// Sunucuyu başlat
	log.Fatal(http.ListenAndServe(cfg.GetAddress(), nil))
}
