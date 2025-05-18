FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/waifu-bot ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/waifu-bot .

CMD [ "./waifu-bot" ]