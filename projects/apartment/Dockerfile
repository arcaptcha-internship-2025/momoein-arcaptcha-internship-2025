FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/main.go

FROM alpine:3 AS stage

RUN apk update && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /main ./

ENTRYPOINT ["./main"]