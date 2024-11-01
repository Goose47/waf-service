FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
ADD go.sum .
COPY . .
RUN go build -o app ./cmd/detection/detection.go

FROM alpine

WORKDIR /app
COPY ./config ./config
COPY --from=builder /build/app /app/app

CMD ["./app"]