FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
ADD go.sum .
COPY . .
RUN go build -o app ./main.go

FROM golang:alpine

WORKDIR /app
COPY . .

CMD ["go",  "test", "tests/..."]