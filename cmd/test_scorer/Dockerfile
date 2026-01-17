# 1. Aşama: Builder
# 1.25 yerine stabil olan 1.23-alpine kullanıyoruz
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Build için git gerekebilir, ekliyoruz
RUN apk add --no-cache git

# Bağımlılıkları önceden çekmek build süresini hızlandırır
COPY go.mod go.sum ./
RUN go mod download

# Tüm kaynak kodunu kopyala
COPY . .

# MİMARİ BURADA KRİTİK: 
# CGO_ENABLED=0: Tamamen statik bir binary oluşturur (Alpine'de sorunsuz çalışması için şart).
# GOOS=linux ve GOARCH=amd64: AMD işlemcin için doğru hedefi belirler.
# Binary ismini ve yolunu netleştiriyoruz.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/scraper_bin ./internal/scraper

# 2. Aşama: Final (Run)
FROM alpine:latest

# Tor ve HTTPS bağlantıları için sertifikalar şart
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Sadece ihtiyacımız olan binary'yi builder'dan alıyoruz
COPY --from=builder /app/scraper_bin ./scraper

# Çalıştırma izni
RUN chmod +x ./scraper

# Go uygulamanın içindeki port ile uyumlu olmalı
EXPOSE 8080

CMD ["./scraper"]