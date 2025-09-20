FROM golang:alpine AS builder
WORKDIR /app
COPY Gateway/go.mod Gateway/go.sum ./
COPY Proto ./Proto
COPY Gateway/ ./
RUN go mod download
RUN go build -o gateway

FROM alpine
COPY --from=builder /app/gateway /usr/bin
EXPOSE 8070
ENTRYPOINT [ "gateway" ]
