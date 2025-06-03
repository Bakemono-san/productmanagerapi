package utils

import (
	"encoding/json"
	"net/http"
	"productmanagerapi/config"
	responseFormatter "productmanagerapi/responseFormatter"

	jwt "github.com/golang-jwt/jwt/v5"
)

var RequestMethodValidator = func(w http.ResponseWriter, r http.Request, method string) bool {
	if r.Method != method {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusMethodNotAllowed, "This method is not allowed", nil))
		return false
	}

	return true
}

var IsValidToken = func(token string) bool {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.SECRET_KEY), nil // Replace with your actual secret key
	})

	if err != nil || !jwtToken.Valid {
		return false
	}
	return true
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if it's swagger or docs route
		if r.URL.Path == "/swagger/index.html" || r.URL.Path == "/docs" || r.URL.Path == "/swagger/doc.json" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized - no token", http.StatusUnauthorized)
			return
		}

		token := cookie.Value

		if token == "" || !IsValidToken(token) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusUnauthorized, "Unauthorized", nil))
			return
		}

		next(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from your frontend (e.g., Next.js running on localhost:3000)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight request (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
