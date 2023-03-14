package services

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (s *userService) Login(name string, password string) (string, error) {
	user, err := s.repo.GetByUsername(name)
	if err != nil {
		return "", err
	}
	// Compare the stored password hash with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        user.ID,
		"email":     user.Email,
		"expiresAt": time.Now().Add(time.Hour * 24).Unix(),
	})

	// tokeString, err := token.SignedString()

	return "", nil
}
