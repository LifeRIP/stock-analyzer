FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN go install github.com/air-verse/air@latest

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
# COPY . .

# Compilar la aplicación
# RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend ./cmd

# Imagen final
# FROM alpine:3.21

# WORKDIR /app

# Copiar el binario compilado
# COPY --from=builder /app/backend .

# Exponer el puerto
EXPOSE 8081

# Ejecutar la aplicación
CMD ["air", "-c", ".air.toml"]