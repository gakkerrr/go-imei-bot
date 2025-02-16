FROM golang:1.23-alpine
WORKDIR /go-imei-bot
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-imei-bot cmd/go-imei-bot/main.go
RUN chmod +x /go-imei-bot
ENV CONFIG_PATH=/go-imei-bot/config/config.yaml 
CMD ["./main"]