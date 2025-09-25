FROM golang:alpine AS builder

# Set Go proxy and other environment variables
ENV GOPROXY=https://proxy.golang.org,https://goproxy.io,https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.org
ENV CGO_ENABLED=0

WORKDIR /app/Stakeholders

COPY Stakeholders/go.mod .
COPY Stakeholders/go.sum .
COPY common /app/common

# Retry download with timeout
RUN for i in 1 2 3; do go mod download && break || sleep 5; done
COPY Stakeholders/. . 

# Retry tidy with timeout  
RUN for i in 1 2 3; do go mod tidy && break || sleep 5; done

RUN go build -o /app/stakeholders .

FROM alpine
COPY --from=builder /app/stakeholders /usr/bin/stakeholders
EXPOSE 8080 50051
ENTRYPOINT ["stakeholders"]
