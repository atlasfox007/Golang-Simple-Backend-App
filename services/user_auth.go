package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/atlasfox007/Golang-Simple-Backend-App/model"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrUsernameAlreadyRegistered = errors.New("Username has already been registered")
	ErrFailedToCreateNewUser     = errors.New("failed to create new user")
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

	tokenString, err := token.SignedString([]byte(os.Getenv("REFRESH_JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *userService) Register(name string, password string, email string) error {
	user, err := s.repo.GetByUsername(name)
	if err != nil {
		return err
	}

	if user != nil {
		return ErrUsernameAlreadyRegistered
	}

	//  Hash User password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password %v", err)
	}
	// Create user
	newUser := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Token:    "",
	}

	err = s.repo.CreateUser(newUser)
	if err != nil {
		return ErrFailedToCreateNewUser
	}
	return nil
}
