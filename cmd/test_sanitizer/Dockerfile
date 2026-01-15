FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /test_sanitizer ./cmd/test_sanitizer

CMD ["/test_sanitizer"]