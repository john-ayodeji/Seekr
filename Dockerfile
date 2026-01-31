# ----------- Build stage -----------
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

# ----------- Run stage -----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/server .

# App port (Docker-side)
EXPOSE 3000

ENV PORT=3000

# Run
CMD ["./server"]