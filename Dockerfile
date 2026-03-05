# Stage de Desarrollo
FROM golang:1.25-alpine AS dev
WORKDIR /app

# Instalamos Air para el Hot Reload
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Ejecutamos air (que leerá el .air.toml)
CMD ["air"]