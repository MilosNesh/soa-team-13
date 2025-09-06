FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o blog

FROM alpine
WORKDIR /app
COPY --from=builder /app/blog /usr/bin
EXPOSE 8081
ENTRYPOINT ["blog"]
