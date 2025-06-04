package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	config "productmanagerapi/config"
	"productmanagerapi/models"
	responseFormatter "productmanagerapi/responseFormatter"
	"productmanagerapi/services"
	"productmanagerapi/utils"
	requestMethodValidator "productmanagerapi/utils"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"

	bcrypt "golang.org/x/crypto/bcrypt"
)

var Login = func(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodPost)
	if !isValidMethod {
		return
	}

	fmt.Println("Processing user login...")

	var userInfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		utils.ResponseWritter(w, http.StatusBadRequest, responseFormatter.FormatResponse(http.StatusBadRequest, "Invalid request body", nil))
		fmt.Println("Invalid request body for login:", err)
		return
	}

	username := userInfo.Username
	password := userInfo.Password

	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		utils.ResponseWritter(w, http.StatusBadRequest, responseFormatter.FormatResponse(http.StatusBadRequest, "Username and password are required", nil))
		return
	}

	tokenString, err, user := services.Login(username, password)

	if err != nil {
		if err.Error() == "invalid username or password" {
			utils.ResponseWritter(w, http.StatusUnauthorized, responseFormatter.FormatResponse(http.StatusUnauthorized, "Invalid username or password", nil))
			fmt.Println("Invalid username or password for user:", username)
			return
		}

		utils.ResponseWritter(w, http.StatusInternalServerError, responseFormatter.FormatResponse(http.StatusInternalServerError, "Error logging in", nil))
		fmt.Println("Error logging in for user:", username, "Error:", err)
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

	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Login successful", map[string]interface{}{
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

var RefreshToken = func(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse and validate existing token
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	// Create a new token with same claims and new expiration
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  claims["user_id"],
		"username": claims["username"],
		"email":    claims["email"],
		"role":     claims["role"],
	})

	tokenString, err := newToken.SignedString([]byte(config.SECRET_KEY))
	if err != nil {
		http.Error(w, "Failed to refresh", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   3600,
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Token refreshed",
	})
}

var Logout = func(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Logged out successfully", nil))
	fmt.Println("User logged out successfully")
}
