# build stage
FROM golang:1.16-alpine3.13 AS builder
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o go-migrate .

# final stage
FROM scratch
COPY --from=builder /app/go-migrate /app/go-migrate
CMD ["/app/go-migrate"]