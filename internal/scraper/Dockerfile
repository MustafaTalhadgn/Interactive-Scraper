# 1. Aşama: Builder
FROM golang:1.25.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Bağımlılıkları kopyala ve indir
COPY go.mod go.sum ./
RUN go mod download

# Tüm projeyi kopyala
COPY . .

# KRİTİK DÜZELTME: 
# cmd/scraper altındaki main.go'yu derliyoruz. 
# Çıktı ismini "cti_scraper_bin" yapalım ki klasör isimleriyle karışmasın.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/cti_scraper_bin ./cmd/scraper/main.go

# 2. Aşama: Final
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Builder'dan gelen binary'yi kopyala
COPY --from=builder /app/cti_scraper_bin ./scraper

# Çalıştırma izni
RUN chmod +x ./scraper

EXPOSE 8080

# Binary'yi çalıştır
CMD ["./scraper"]