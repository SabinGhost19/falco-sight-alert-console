# Etapa de build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Descărcăm dependențele
COPY go.mod go.sum ./
RUN go mod download

# Copiem codul sursă
COPY . .

# Compilăm aplicația (fără CGO, ideal pentru containere lightweight/scratch)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o falcosight-backend ./cmd/api

# Etapa finală - lightweight
FROM alpine:3.19

WORKDIR /app

# Instalăm certificatele root necesare pentru cereri HTTPS externe (ex. Talon)
RUN apk --no-cache add ca-certificates

# Copiem doar binarul compilat
COPY --from=builder /app/falcosight-backend .

EXPOSE 8080

CMD ["./falcosight-backend"]
