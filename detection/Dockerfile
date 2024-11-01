FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
ADD go.sum .
COPY . .
RUN go build -o app ./cmd/waf/waf.go

FROM alpine

WORKDIR /app
COPY --from=builder /build/app /app/app

CMD ["./app"]