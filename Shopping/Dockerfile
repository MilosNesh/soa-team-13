FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o shopping

FROM alpine
WORKDIR /app
COPY --from=builder /app/shopping /usr/bin
EXPOSE 8082
ENTRYPOINT ["shopping"]
