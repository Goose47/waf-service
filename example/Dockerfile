FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
ADD go.sum .
COPY . .
RUN go build -o app main.go

FROM alpine

WORKDIR /app
COPY --from=builder /build/app /app/app
COPY ./assets /app/assets

CMD ["./app"]