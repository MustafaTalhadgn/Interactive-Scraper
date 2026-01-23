
FROM golang:1.25.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/cti_scraper_bin ./cmd/scraper/main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/cti_scraper_bin ./scraper


RUN chmod +x ./scraper

EXPOSE 8080


CMD ["./scraper"]