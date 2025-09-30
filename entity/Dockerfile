FROM golang:alpine AS builder
WORKDIR /app
COPY go.* .
COPY communication ./communication
COPY entity ./entity
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main entity/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
