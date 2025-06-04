package services

import (
	"fmt"
	"productmanagerapi/config"
	"productmanagerapi/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var Login = func(username, password string) (string, error, models.User) {
	// Fetch user from database
	var user models.User
	result := config.Db.Where("username = ? ", username).First(&user)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			fmt.Println("Invalid username or password for user:", username)
			return "", fmt.Errorf("invalid username or password"), models.User{}
		}
		fmt.Println("Error fetching user:", result.Error)
		return "", fmt.Errorf("error fetching user: %v", result.Error), models.User{}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		fmt.Println("Invalid username or password for user:", username)
		return "", fmt.Errorf("invalid username or password"), models.User{}
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
		return "", err, models.User{}
	}

	return tokenString, nil, user
}
