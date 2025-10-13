package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"gateway.com/dto"
)

type GatewayHandler struct {
}

func (handler *GatewayHandler) HandleAccount(writer http.ResponseWriter, req *http.Request) {
	var tokenData *dto.TokenData = parseToken(req)
	path := strings.TrimPrefix(req.URL.Path, "/accounts/")

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
	log.Printf("DEBUG Gateway: === INCOMING REQUEST ===")
	log.Printf("DEBUG Gateway: Method: %s", req.Method)
	log.Printf("DEBUG Gateway: Path: %s", req.URL.Path)
	log.Printf("DEBUG Gateway: Authorization: %s", req.Header.Get("Authorization"))
	log.Printf("DEBUG Gateway: HandleBlog received request for path: %s", req.URL.Path)
	var tokenData *dto.TokenData = parseToken(req)
	path := strings.TrimPrefix(req.URL.Path, "/blogs/")

	log.Printf("DEBUG Gateway: Trimmed path: %s", path)

	//Prosljedjivanje zahtjeva
	targetURL, err := url.Parse("http://blog_service:8081")
	if err != nil {
		http.Error(writer, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	log.Printf("DEBUG Gateway: Target URL for proxy: %s", targetURL.String())

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = "/blogs/" + path
		req.Host = "blog_service:8081"

		log.Printf("DEBUG Gateway: Proxying to new path: %s on host: %s", req.URL.Path, req.Host)

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
		log.Println("DEBUG: prazan header")
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

func (handler *GatewayHandler) HandleShopping(writer http.ResponseWriter, req *http.Request) {
	var tokenData *dto.TokenData = parseToken(req)
	path := strings.TrimPrefix(req.URL.Path, "/shopping/")

	if tokenData == nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Prosljedjivanje zahtjeva
	targetURL, err := url.Parse("http://shopping_service:8082")
	if err != nil {
		http.Error(writer, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = "/shopping/" + path
		req.Host = "shopping_service:8082"

		if tokenData != nil {
			req.Header.Set("X-Account-Username", tokenData.Username)
			req.Header.Set("X-Account-Role", tokenData.Role)
			req.Header.Set("X-Account-Id", tokenData.Id)
		}
	}

	proxy.ServeHTTP(writer, req)
}

func (handler *GatewayHandler) HandleStaticFiles(writer http.ResponseWriter, req *http.Request) {
	targetURL, err := url.Parse("http://blog_service:8081")
	if err != nil {
		http.Error(writer, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		req.URL.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme
		req.Host = targetURL.Host

		req.Body = nil
		if req.GetBody != nil {
			req.GetBody = func() (closer io.ReadCloser, e error) { return nil, nil }
		}

	}

	proxy.ServeHTTP(writer, req)
}

func (handler *GatewayHandler) HandleFollow(writer http.ResponseWriter, req *http.Request) {
	var token = parseToken(req)
	path := strings.TrimPrefix(req.URL.Path, "/follow/")

	if token == nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// role, _, _ := parseToken(token) // ako treba rola za neku provjeru otkomentarisati, takodje username i id

	//Prosljedjivanje zahtjeva
	targetURL, err := url.Parse("http://followers_service:8083")
	if err != nil {
		http.Error(writer, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = "/follow/" + path
		req.Host = "followers_service:8083"
	}

	proxy.ServeHTTP(writer, req)
}
