FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . . 
RUN go build -o stakeholders

FROM alpine
COPY --from=builder /app/stakeholders /usr/bin
EXPOSE 8080 50051
ENTRYPOINT ["stakeholders"]
