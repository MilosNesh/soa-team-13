package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
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

func (handler *GatewayHandler) HandleTour(writer http.ResponseWriter, req *http.Request) {
	var token = getToken(req)
	path := strings.TrimPrefix(req.URL.Path, "/tours/")

	// Proveri token, osim ako je GET metod
	if token == "" && !(path == "" && req.Method == "GET") {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	id, role := "", ""
	if token != "" {
		role, _, id = parseToken(token)
	}

	if req.Method == "POST" && path == "" && role != "guide" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	if req.Method == "POST" && path == "" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(writer, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer req.Body.Close()

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(writer, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		if _, exists := data["authorId"]; exists {
			data["authorId"] = id
		}

		updatedBody, err := json.Marshal(data)
		if err != nil {
			http.Error(writer, "Failed to re-encode request body", http.StatusInternalServerError)
			return
		}

		req.Body = io.NopCloser(bytes.NewReader(updatedBody))
		req.ContentLength = int64(len(updatedBody))
		req.Header.Set("Content-Length", strconv.Itoa(len(updatedBody)))
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
	}

	// Poslati zahtev (s izmenjenim telom ako je POST)
	proxy.ServeHTTP(writer, req)
}

func parseToken(token string) (string, string, string) {
	secretKey := []byte("sekret_key_12#4")

	parsedToken, err := jwt.ParseWithClaims(token, &dto.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Proveri da li je algoritam ispravan
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
