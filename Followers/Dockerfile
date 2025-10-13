FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o followers

FROM golang:alpine
WORKDIR /app
COPY --from=builder /app/followers /usr/bin
EXPOSE 8082
ENTRYPOINT ["followers"]