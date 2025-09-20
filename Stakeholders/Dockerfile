FROM golang:alpine AS builder
WORKDIR /app

COPY Stakeholders/go.mod Stakeholders/go.sum ./
COPY Proto ./Proto
COPY Stakeholders/ ./

RUN go mod download
RUN go build -o stakeholders

FROM alpine
COPY --from=builder /app/stakeholders /usr/bin
EXPOSE 8080 50051
ENTRYPOINT ["stakeholders"]
