# 🔗 URL Kısaltma Servisi

Modern ve hızlı bir URL kısaltma servisi. Go (Golang) ile geliştirilmiş, basit ve etkili bir web uygulaması.

## 📋 İçindekiler

- [Özellikler](#-özellikler)
- [Teknolojiler](#-teknolojiler)
- [Proje Yapısı](#-proje-yapısı)
- [Kurulum](#-kurulum)
- [Kullanım](#-kullanım)
- [API Endpoints](#-api-endpoints)
- [Environment Variables](#-environment-variables)
- [Örnekler](#-örnekler)
- [Katkıda Bulunma](#-katkıda-bulunma)

## ✨ Özellikler

- 🚀 **Hızlı URL Kısaltma**: Uzun URL'leri kısa kodlara dönüştürme
- 📊 **İstatistik Takibi**: Tıklanma sayısı ve detaylı istatistikler
- 🌐 **Web Arayüzü**: Kullanıcı dostu HTML arayüzü
- 🔄 **Otomatik Yönlendirme**: Kısa URL'lerden orijinal URL'e yönlendirme
- ⚙️ **Environment Variables**: Esnek konfigürasyon desteği
- 💾 **In-Memory Database**: Hızlı veri saklama (geliştirme için)
- 💾 **SQLite Database**: SQLite kullanarak kalıcı veri saklama
- 💾 **Bolt Database**: BBolt kullanarak kalıcı veri saklama 

## 🛠 Teknolojiler

- **Backend**: Go (Golang) 1.24.5
- **HTTP Server**: Go net/http (standart kütüphane)
- **Frontend**: HTML5, Tailwind, Vanilla JavaScript
- **Modül Sistemi**: Go Modules
- **Veritabanı**: In-Memory, SQLite, Bolt

## 📁 Proje Yapısı

```
urlShortener/
├── 📄 main.go             # Ana uygulama dosyası
├── 📄 go.mod              # Go modül tanımları
├── 📄 .env                # Environment değişkenleri
├── 📁 config/             # Konfigürasyon yönetimi
│   └── config.go
├── 📁 handlers/           # HTTP route handlers
│   └── url_handler.go
├── 📁 storage/             # Veri modelleri
│   └── bolt.go
│   └── url.go
├── 📁 utils/              # Yardımcı fonksiyonlar
│   ├── env.go
│   └── shortener.go
└── 📁 pages/              # Frontend dosyaları
    └── index.html
```

## 🚀 Kurulum

### Gereksinimler

- Go 1.20 veya üzeri

### Adımlar

1. **Projeyi klonlayın veya indirin:**
   ```bash
   git clone <repository-url>
   cd urlShortener
   ```

2. **Bağımlılıkları kontrol edin:**
   ```bash
   go mod tidy
   ```

3. **Uygulamayı çalıştırın:**
   ```bash
   go run main.go
   ```

4. **Tarayıcınızda açın:**
   ```
   http://localhost:8080
   ```

## 💻 Kullanım

### 1. Web Arayüzü ile Kullanım

1. Tarayıcınızda `http://localhost:8080` adresine gidin
2. URL input alanına kısaltmak istediğiniz URL'i girin
3. "Kısalt" butonuna tıklayın
4. Oluşturulan kısa URL'i kopyalayın ve paylaşın

### 2. API ile Kullanım

#### URL Kısaltma:
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.example.com/very/long/url"}'
```

**Yanıt:**
```json
{
  "original": "https://www.example.com/very/long/url",
  "short_url": "http://localhost:8080/s/abc123",
  "short_code": "abc123"
}
```

## 🔗 API Endpoints

| Method | Endpoint | Açıklama |
|--------|----------|----------|
| `GET` | `/` | Ana sayfa (HTML arayüzü) |
| `POST` | `/shorten` | URL kısaltma |
| `GET` | `/s/{code}` | Kısa URL'den yönlendirme |
| `GET` | `/stats/{code}` | URL istatistikleri |

### Detaylı API Dokumentasyonu

#### POST /shorten
URL kısaltma işlemi yapar.

**Request Body:**
```json
{
  "url": "https://example.com"
}
```

**Response:**
```json
{
  "original": "https://example.com",
  "short_url": "http://localhost:8080/s/abc123",
  "short_code": "abc123"
}
```

#### GET /stats/{code}
Belirli bir kısa URL'in istatistiklerini getirir.

**Response:**
```json
{
  "id": 1,
  "original": "https://example.com",
  "short": "abc123",
  "created_at": "2025-08-01T10:30:00Z",
  "clicks": 5
}
```

## ⚙️ Environment Variables

Uygulamayı farklı konfigürasyonlarla çalıştırabilirsiniz:

| Değişken | Açıklama | Varsayılan |
|----------|----------|------------|
| `PORT` | Server portu                     | `8080` |
| `HOST` | Server host adresi               | `localhost` |
| `SSL`  | Server HTTPS mi tercih etsin     | `localhost` |
| `DEBUG`  | Debug modu                     | `false` |
| `SHORTENER_LENGTH`  | Kısaltma uzunluğu   | `6` |
| `DATABASE` | Veritabanı tipi              | `memory` |

## 📝 Örnekler

### 1. Farklı Port ile Çalıştırma

**PowerShell:**
```powershell
$env:PORT="3000"; go run main.go
```

**Command Prompt:**
```cmd
set PORT=3000 && go run main.go
```

**Linux/Mac:**
```bash
PORT=3000 go run main.go
```

### 2. Debug Modu ile Çalıştırma

```powershell
$env:DEBUG="true"; go run main.go
```

### 3. Farklı Host ve Port

```powershell
$env:HOST="0.0.0.0"; $env:PORT="9000"; go run main.go
```

### 4. Tüm Ayarlar ile

```powershell
$env:PORT="5000"; $env:DEBUG="true"; $env:HOST="127.0.0.1"; go run main.go
```

## 🔧 Geliştirme

### Kod Yapısı

- **main.go**: Uygulama başlangıç noktası ve route tanımları
- **config/**: Konfigürasyon yönetimi ve environment variables
- **handlers/**: HTTP request handlers
- **storage/**: Veri yapıları ve database işlemleri
- **utils/**: Yardımcı fonksiyonlar (URL kısaltma, environment handling)
- **pages/**: Frontend HTML dosyaları

### Yeni Özellik Ekleme

1. **Yeni Route Eklemek:**
   ```go
   // main.go içinde
   http.HandleFunc("/new-route", handlers.NewHandler)
   ```

2. **Yeni Handler Oluşturmak:**
   ```go
   // handlers/url_handler.go içinde
   func NewHandler(w http.ResponseWriter, r *http.Request) {
       // Handler logic
   }
   ```

## 🤝 Katkıda Bulunma

1. Repository'yi fork edin
2. Yeni bir branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'Add amazing feature'`)
4. Branch'inizi push edin (`git push origin feature/amazing-feature`)
5. Pull Request oluşturun

## 📄 Lisans

Bu proje [MIT License](LICENSE) altında lisanslanmıştır.

## 🐛 Sorun Bildirme

Bir sorun bulduysanız veya öneriniz varsa, lütfen [Issues](../../issues) bölümünden bildirin.

## 🎯 Gelecek Planları

- [x] SQLite veritabanı desteği
- [x] BBolt veritabanı desteği
- [ ] PostgreSQL desteği
- [ ] Redis cache entegrasyonu
- [ ] Custom short URL kodları
- [ ] URL expiration (son kullanma tarihi)
- [ ] Bulk URL shortening
- [ ] QR code üretimi
- [ ] Analytics dashboard
- [ ] Rate limiting
- [ ] API authentication