FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

FROM alpine:3.14 AS final

COPY --from=builder /build/server /bin/server

ENTRYPOINT ["/bin/server", "filmoteka"]