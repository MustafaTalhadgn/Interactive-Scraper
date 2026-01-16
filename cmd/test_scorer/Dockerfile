FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /test_scorer ./cmd/test_scorer

CMD ["/test_scorer"]