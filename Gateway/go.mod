module gateway.com

go 1.24.4

require (
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	google.golang.org/grpc v1.75.1
	grpc/proto v0.0.1
)

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace grpc/proto => ./Proto
