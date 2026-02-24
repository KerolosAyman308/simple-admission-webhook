FROM golang:1.26-alpine3.23 as builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o /build/app ./cmd

FROM alpine:3.23 AS final
WORKDIR /app

COPY --from=builder /build/app ./app

EXPOSE 8000

CMD ["./app"]