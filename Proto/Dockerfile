FROM golang:1.18-alpine AS builder

# Postavljanje radnog direktorijuma
WORKDIR /app

# Kopiranje go.mod i go.sum iz root foldera
COPY ../go.mod ../go.sum ./

# Preuzimanje zavisnosti
RUN go mod tidy
RUN go mod download

# Kopiranje 'proto' direktorijuma iz root foldera
COPY ../proto ./proto

# Instalacija potrebnih alata
RUN apk update && apk add --no-cache \
    bash \
    protobuf \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.0 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

# Pokretanje protoc za generisanje Go koda
RUN protoc --go_out=. --go-grpc_out=. --proto_path=./proto ./proto/*.proto

# Build proto-generator
RUN go build -o proto-generator

# Drugi deo za pokretanje aplikacije
FROM alpine
WORKDIR /root/

# Kopiraj izgrađeni binarni fajl iz prethodnog stanja
COPY --from=builder /app/proto-generator /usr/bin/

EXPOSE 50051

ENTRYPOINT [ "proto-generator" ]
