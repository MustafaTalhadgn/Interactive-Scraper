FROM golang:1.23-alpine AS builder

WORKDIR /app


RUN apk add --no-cache git


COPY go.mod go.sum ./
RUN go mod download


COPY . .


WORKDIR /app/internal/scraper
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scraper .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/


COPY --from=builder /app/internal/scraper/scraper .


EXPOSE 8080


CMD ["./scraper"]