# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar los archivos del módulo primero para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o resource-service ./cmd/api/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /app/resource-service .


# Puerto expuesto (debe coincidir con tu configuración HTTP)
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./resource-service"]