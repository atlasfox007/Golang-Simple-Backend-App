package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/atlasfox007/Golang-Simple-Backend-App/model"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrUsernameAlreadyRegistered = errors.New("username has already been registered")
	ErrFailedToCreateNewUser     = errors.New("failed to create new user")
	ErrFailedToUpdateJWTToken    = errors.New("failed to update jwt token in database")
	ErrFailedToSignJWTToken      = errors.New("failed to sign jwt token")
)

func (s *userService) Login(name string, password string) (string, error) {
	user, err := s.repo.GetByUsername(name)

	if err != nil {
		return "", errors.New("failed to get the user by username")
	}
	// Compare the stored password hash with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  user.Name,
		"expiresAt": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Load JWT Secret from .env files
	err = godotenv.Load()
	if err != nil {
		return "", errors.New("failed to load environment files")
	}

	val, found := os.LookupEnv("REFRESH_JWT_SECRET")
	if !found {
		return "", errors.New("failed to get the refresh jwt secret")
	}

	tokenString, err := token.SignedString([]byte(val))
	if err != nil {
		return "", ErrFailedToSignJWTToken
	}

	// Put the token inside the user database
	user.Token = tokenString

	err = s.repo.UpdateUser(user)
	if err != nil {
		return "", ErrFailedToUpdateJWTToken
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
