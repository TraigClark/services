# Dockerfile in services/hello-world/src

# Build stage
FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/cmd/main .
COPY --from=build /app/go.mod .
COPY --from=build /app/go.sum .

# Copy the configuration file
COPY --from=build /app/cmd/config/config.json ./config/config.json

# EXPOSE PORTS
EXPOSE 1883
EXPOSE 1502

CMD ["./main"]