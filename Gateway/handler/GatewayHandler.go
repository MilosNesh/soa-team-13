package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"gateway.com/dto"
	"github.com/golang-jwt/jwt"
)

type GatewayHandler struct{}

func (handler *GatewayHandler) HandleAccount(writer http.ResponseWriter, req *http.Request) {
	var role = ""
	var token = getToken(req)
	path := strings.TrimPrefix(req.URL.Path, "/accounts/")

	if token == "" && !(path == "" && req.Method == "POST") && !(path == "login" && req.Method == "POST") {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if token != "" {
		role, _, _ = parseToken(token)
	}

	if path == "" && req.Method == "GET" && role != "admin" {
		writer.WriteHeader(http.StatusUnauthorized)
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
	}

	proxy.ServeHTTP(writer, req)
}

func (handler *GatewayHandler) HandleBlog(writer http.ResponseWriter, req *http.Request) {
	var token = getToken(req)
	path := strings.TrimPrefix(req.URL.Path, "/blogs/")

	if token == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// role, _, _ := parseToken(token) // ako treba rola za neku provjeru otkomentarisati, takodje username i id

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
	}

	proxy.ServeHTTP(writer, req)
}

func parseToken(token string) (string, string, string) {
	secretKey := []byte("sekret_key_12#4")

	parsedToken, err := jwt.ParseWithClaims(token, &dto.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("nevalidan signing metod: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		log.Fatalf("Greška pri parsiranju tokena: %v", err)
	}

	role, username, id := "", "", ""
	if claims, ok := parsedToken.Claims.(*dto.Claims); ok && parsedToken.Valid {
		role = claims.Role
		username = claims.Username
		id = claims.Subject
	}
	return role, username, id
}

func getToken(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")

	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func (handler *GatewayHandler) HandleFollow(writer http.ResponseWriter, req *http.Request) {
	var token = getToken(req)
	path := strings.TrimPrefix(req.URL.Path, "/follow/")

	if token == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// role, _, _ := parseToken(token) // ako treba rola za neku provjeru otkomentarisati, takodje username i id

	//Prosljedjivanje zahtjeva
	targetURL, err := url.Parse("http://followers_service:8082")
	if err != nil {
		http.Error(writer, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = "/follow/" + path
		req.Host = "followers_service:8082"
	}

	proxy.ServeHTTP(writer, req)
}
