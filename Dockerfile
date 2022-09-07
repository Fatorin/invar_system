FROM golang:1.18 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:3.15.0
#FROM scratch
WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["./main"]