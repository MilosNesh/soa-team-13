package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"grpc/proto"

	"gateway.com/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type GatewayHandler struct {
	grpcClient proto.StakeholdersServiceClient
}

func (handler *GatewayHandler) HandleAccount(writer http.ResponseWriter, req *http.Request) {
	var tokenData *dto.TokenData = parseToken(req)
	path := strings.TrimPrefix(req.URL.Path, "/accounts/")

	if tokenData == nil && !(path == "" && req.Method == "POST") && !(path == "login" && req.Method == "POST") {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if path == "login" && req.Method == "POST" {
		handler.handleLoginWithGRPC(writer, req)
		return
	}

	if path == "" && req.Method == "POST" {
		handler.handleRegisterWithGRPC(writer, req)
		return
	}

	//Prosljedjivanje zahtjeva
	targetURL, err := url.Parse("http://stakeholders_service:8080")
	if err != nil {
		http.Error(writer, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = "/accounts/" + path
		req.Host = "stakeholders_service:8080"

		if tokenData != nil {
			req.Header.Set("X-Account-Username", tokenData.Username)
			req.Header.Set("X-Account-Role", tokenData.Role)
			req.Header.Set("X-Account-Id", tokenData.Id)
		}
	}

	proxy.ServeHTTP(writer, req)
}

func (handler *GatewayHandler) HandleBlog(writer http.ResponseWriter, req *http.Request) {
	var tokenData *dto.TokenData = parseToken(req)
	path := strings.TrimPrefix(req.URL.Path, "/blogs/")

	if tokenData == nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Prosljedjivanje zahtjeva
	targetURL, err := url.Parse("http://blog_service:8081")
	if err != nil {
		http.Error(writer, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = "/blogs/" + path
		req.Host = "blog_service:8081"

		if tokenData != nil {
			req.Header.Set("X-Account-Username", tokenData.Username)
			req.Header.Set("X-Account-Role", tokenData.Role)
			req.Header.Set("X-Account-Id", tokenData.Id)
		}
	}

	proxy.ServeHTTP(writer, req)
}

func (handler *GatewayHandler) HandleTour(writer http.ResponseWriter, req *http.Request) {
	var tokenData *dto.TokenData = parseToken(req)
	path := strings.TrimPrefix(req.URL.Path, "/tours/")

	// Proveri token, osim ako je GET metod
	if tokenData == nil && !(path == "" && req.Method == "GET") {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	targetURL, err := url.Parse("http://tours_service:8084")
	if err != nil {
		http.Error(writer, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = "/tours/" + path
		req.Host = "tours_service:8084"

		if tokenData != nil {
			req.Header.Set("X-Account-Username", tokenData.Username)
			req.Header.Set("X-Account-Role", tokenData.Role)
			req.Header.Set("X-Account-Id", tokenData.Id)
		}
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		for key := range resp.Header {
			if strings.HasPrefix(strings.ToLower(key), "access-control-") {
				delete(resp.Header, key)
			}
		}
		return nil
	}

	proxy.ServeHTTP(writer, req)
}

func parseToken(r *http.Request) *dto.TokenData {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil
	}

	req, _ := http.NewRequest("GET", "http://stakeholders_service:8080/accounts/parsetoken", nil)
	req.Header.Set("Authorization", authHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}
	defer resp.Body.Close()

	var tokenData *dto.TokenData

	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		return nil
	}
	return tokenData
}

func (handler *GatewayHandler) ConnectToGRPCServer() {
	conn, err := grpc.Dial("stakeholders_service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ne može se povezati na gRPC server: %v", err)
	}
	handler.grpcClient = proto.NewStakeholdersServiceClient(conn)
}

func (handler *GatewayHandler) handleLoginWithGRPC(writer http.ResponseWriter, req *http.Request) {
	println("Login sa grpc")
	var loginReq proto.LoginRequest
	err := json.NewDecoder(req.Body).Decode(&loginReq)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	loginResp, err := handler.grpcClient.Login(context.Background(), &loginReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			http.Error(writer, st.Message(), http.StatusBadRequest)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf(`{"token":"%s"}`, loginResp.Token)))
}

func (handler *GatewayHandler) handleRegisterWithGRPC(writer http.ResponseWriter, req *http.Request) {
	println("Register sa grpc")
	var registerReq proto.RegisterAccountRequest
	err := json.NewDecoder(req.Body).Decode(&registerReq)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	regResp, err := handler.grpcClient.RegisterAccount(context.Background(), &registerReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			http.Error(writer, st.Message(), http.StatusBadRequest)
			return
		}
	}
	println(regResp)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(`{"message": "Uspješno registracija"}`))
}
