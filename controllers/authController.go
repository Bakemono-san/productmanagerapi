package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	config "productmanagerapi/config"
	"productmanagerapi/models"
	responseFormatter "productmanagerapi/responseFormatter"
	requestMethodValidator "productmanagerapi/utils"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"

	bcrypt "golang.org/x/crypto/bcrypt"
)

var Login = func(w http.ResponseWriter, r *http.Request) {
	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodPost)
	if !isValidMethod {
		return
	}

	fmt.Println("Processing user login...")
	w.Header().Set("Content-Type", "application/json")

	var userInfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Invalid request body", nil))
		fmt.Println("Invalid request body for login:", err)
		return

	}
	username := userInfo.Username
	password := userInfo.Password

	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Username and password are required", nil))
		return
	}

	var user models.User
	result := config.Db.Where("username = ? ", username).First(&user)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusUnauthorized, "Invalid username or password", nil))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching user", nil))
		}
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusUnauthorized, "Invalid username or password", nil))
		fmt.Println("Invalid username or password for user:", username)
		return
	}

	// generate a token for the user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})

	tokenString, err := token.SignedString([]byte(config.SECRET_KEY))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error generating token", nil))
		fmt.Println("Error generating token for user:", username, "Error:", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   3600, // 1 hour
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Login successful", map[string]interface{}{
		"user_id":      user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"role":         user.Role,
		"access_token": tokenString,
	}))
	fmt.Println("User login successful")
}

var Register = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodPost)
	if !isValidMethod {
		return
	}

	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
	fmt.Println("Processing user registration...")
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Invalid request body", nil))
		return
	}

	if user.Username == "" || user.Password == "" || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Username, password, and email are required", nil))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error hashing password", nil))
		fmt.Println("Error hashing password:", err)
		return
	}
	user.Password = string(hashedPassword)

	if err := config.Db.Create(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error creating user", nil))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusCreated, "User registered successfully", user))
	fmt.Println("User registration successful")
}

var HomeController = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Welcome to the Product Manager API", nil))
}
