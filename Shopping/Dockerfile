FROM golang:alpine AS builder

# Set Go proxy and other environment variables
ENV GOPROXY=https://proxy.golang.org,https://goproxy.io,https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.org
ENV CGO_ENABLED=0

WORKDIR /app/Shopping
COPY Shopping/go.mod .
COPY Shopping/go.sum .
COPY common /app/common

# Retry download with timeout
RUN for i in 1 2 3; do go mod download && break || sleep 5; done
COPY Shopping/. .

# Retry tidy with timeout  
RUN for i in 1 2 3; do go mod tidy && break || sleep 5; done

RUN go build -o /app/shopping .

FROM alpine
WORKDIR /app
COPY --from=builder /app/shopping /usr/bin/shopping
EXPOSE 8082
ENTRYPOINT ["shopping"]
